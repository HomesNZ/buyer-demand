package address

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
	"strings"
)

type Repo interface {
	Populate(ctx context.Context, buyerDemands entity.BuyerDemands, needToDeleteTodayData bool) error
	LatestStats(ctx context.Context, suburbID, bedroom, bathroom null.Int, propertyType null.String) (entity.BuyerDemand, error)
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

const latestStatsQuery = `
	SELECT
	    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY current_median_days_to_sell) AS current_median_days_to_sell,
	    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY previous_median_days_to_sell) AS previous_median_days_to_sell,
	    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY current_median_sale_price) AS current_median_sale_price,
	    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY previous_median_sale_price) AS previous_median_sale_price,
	    SUM(num_for_sale_properties) AS num_for_sale_properties,
	    SUM(current_num_for_sale_properties) AS current_num_for_sale_properties,
	    SUM(previous_num_for_sale_properties) AS previous_num_for_sale_properties,
	    MAX(created_at) AS created_at
	FROM homes_data_export.buyer_demand
	WHERE FALSE %s;
`

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

func (r *repo) LatestStats(ctx context.Context, suburbID, bedroom, bathroom null.Int, propertyType null.String) (entity.BuyerDemand, error) {
	whereArray, args := generateWhereClause(suburbID, bedroom, bathroom, propertyType)

	whereClause := fmt.Sprintf(" OR (%s)", strings.Join(whereArray, " AND "))
	query := fmt.Sprintf(latestStatsQuery, whereClause)
	row := r.db.QueryRow(ctx, query, args...)

	bd := entity.BuyerDemand{}
	err := row.Scan(
		&bd.CurrentRangeMedianDaysToSell,
		&bd.PreviousRangeMedianDaysToSell,
		&bd.CurrentRangeMedianSalePrice,
		&bd.PreviousRangeMedianSalePrice,
		&bd.NumOfForSaleProperties,
		&bd.CurrentRangeNumOfForSaleProperties,
		&bd.PreviousRangeNumOfForSaleProperties,
		&bd.CreatedAt,
	)

	return bd, errors.Wrap(err, "rows.Scan")
}
