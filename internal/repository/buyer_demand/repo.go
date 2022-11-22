package address

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
	"strings"
	"time"
)

type Repo interface {
	Populate(ctx context.Context, buyerDemands entity.BuyerDemands, needToDeleteTodayData bool) error
	LatestStats(ctx context.Context, suburbID, bedroom, bathroom null.Int, propertyType null.String) (entity.BuyerDemands, error)
}

func New(db *pgxpool.Pool) (Repo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &repo{db}, nil
}

type repo struct {
	db *pgxpool.Pool
}

const deleteTodayDataQuery = `
	DELETE FROM homes_data_export.buyer_demand
	WHERE created_at > CURRENT_DATE;
`

const populateQuery = `
	INSERT INTO homes_data_export.buyer_demand (
	    suburb_id, 
	    num_bedrooms, 
	    num_bathrooms, 
	    property_type, 
        current_median_days_to_sell,
        previous_median_days_to_sell,
        current_median_sale_price,
        previous_median_sale_price,
        num_for_sale_properties,
        current_num_for_sale_properties,
        previous_num_for_sale_properties,
	    created_at
	)
	VALUES (
	    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now()
	);
`

func (r *repo) Populate(ctx context.Context, buyerDemands entity.BuyerDemands, needToDeleteTodayData bool) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return errors.Wrap(err, "tx.Rollback")
		}
		return errors.Wrap(err, "db.Begin")
	}

	if needToDeleteTodayData {
		_, err = tx.Exec(ctx, deleteTodayDataQuery)
		if err != nil {
			return errors.Wrap(err, "tx.Exec delete")
		}
	}

	for _, bd := range buyerDemands {
		_, err = tx.Exec(
			ctx,
			populateQuery,
			bd.SuburbID,
			bd.NumBedrooms,
			bd.NumBathrooms,
			bd.PropertyType,
			bd.CurrentRangeMedianDaysToSell,
			bd.PreviousRangeMedianDaysToSell,
			bd.CurrentRangeMedianSalePrice,
			bd.PreviousRangeMedianSalePrice,
			bd.NumOfForSaleProperties,
			bd.CurrentRangeNumOfForSaleProperties,
			bd.PreviousRangeNumOfForSaleProperties)
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return errors.Wrap(err, "tx.Rollback")
			}
			return errors.Wrap(err, "tx.Exec")
		}
	}

	err = tx.Commit(ctx)
	return errors.Wrap(err, "tx.Commit")
}

const (
	latestStatsQuery = `
		SELECT
		    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY median_days_to_sell) AS median_days_to_sell,
		    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY median_sale_price) AS median_sale_price,
		    SUM(num_for_sale_properties) AS num_for_sale_properties
		FROM homes_data_export.buyer_demand
		WHERE FALSE %s;
	`
	latestTwoCreatedAtQuery = `
		SELECT DISTINCT created_at 
		FROM homes_data_export.buyer_demand
		WHERE FALSE %s 
		ORDER BY created_at DESC 
		LIMIT 2;
	`
)

func generateWhereClause(suburbID, bedroom, bathroom null.Int, propertyType null.String) ([]string, []interface{}) {
	var whereArray []string
	var values []interface{}
	index := 0

	if !suburbID.IsZero() {
		index++
		whereArray = append(whereArray, fmt.Sprintf("suburb_id = $%d", index))
		values = append(values, suburbID.ValueOrZero())
	}

	if !bedroom.IsZero() {
		index++
		whereArray = append(whereArray, fmt.Sprintf("num_bedrooms = $%d", index))
		values = append(values, bedroom.ValueOrZero())
	}

	if !bathroom.IsZero() {
		index++
		whereArray = append(whereArray, fmt.Sprintf("num_bathrooms = $%d", index))
		values = append(values, bathroom.ValueOrZero())
	}

	if !propertyType.IsZero() {
		index++
		whereArray = append(whereArray, fmt.Sprintf("property_type = $%d", index))
		values = append(values, propertyType.ValueOrZero())
	}

	return whereArray, values
}

func extendWhereClause(whereArray []string, args []interface{}, createdDate time.Time) ([]string, []interface{}) {
	whereArray = append(whereArray, fmt.Sprintf("created_at = $%d", len(args)+1))
	args = append(args, createdDate)
	return whereArray, args
}

func (r *repo) LatestStats(ctx context.Context, suburbID, bedroom, bathroom null.Int, propertyType null.String) (entity.BuyerDemands, error) {
	resp := entity.BuyerDemands{}
	whereArray, args := generateWhereClause(suburbID, bedroom, bathroom, propertyType)

	createdDatesWhereClause := fmt.Sprintf(" OR (%s)", strings.Join(whereArray, " AND "))
	createdDatesQuery := fmt.Sprintf(latestTwoCreatedAtQuery, createdDatesWhereClause)
	createdDatesRows, err := r.db.Query(ctx, createdDatesQuery, args...)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}

	var createdDates []time.Time
	defer createdDatesRows.Close()
	for createdDatesRows.Next() {
		v := time.Time{}
		err := createdDatesRows.Scan(&v)
		if err != nil {
			return nil, errors.Wrap(err, "createdDatesRows.Scan")
		}
		createdDates = append(createdDates, v)
	}

	for _, d := range createdDates {
		extendedWhereArray, extendedArgs := extendWhereClause(whereArray, args, d)
		whereClause := fmt.Sprintf(" OR (%s)", strings.Join(extendedWhereArray, " AND "))
		query := fmt.Sprintf(latestStatsQuery, whereClause)
		row := r.db.QueryRow(ctx, query, extendedArgs...)

		v := entity.BuyerDemand{}
		err := row.Scan(
			&v.CurrentRangeMedianDaysToSell,
			&v.CurrentRangeMedianSalePrice,
			&v.NumOfForSaleProperties,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		v.CreatedAt = null.TimeFrom(d)
		resp = append(resp, v)
	}
	return resp, nil
}
