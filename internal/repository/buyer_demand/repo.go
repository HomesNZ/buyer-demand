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
	    median_days_to_sell, 
	    median_sale_price, 
	    num_for_sale_properties, 
	    created_at
	)
	VALUES (
	    $1, $2, $3, $4, $5, $6, $7, now()
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
			bd.MedianDaysToSell,
			bd.MedianSalePrice,
			bd.NumOfForSaleProperties)
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
		median_days_to_sell,
		median_sale_price,
		num_for_sale_properties,
		created_at
	FROM homes_data_export.buyer_demand
	WHERE FALSE %s
	ORDER BY created_at DESC
	LIMIT 2;
`

func generateWhereClause(suburbID, bedroom, bathroom null.Int, propertyType null.String) (string, []interface{}) {
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

	where := fmt.Sprintf(" OR (%s)", strings.Join(whereArray, " AND "))
	return fmt.Sprintf(latestStatsQuery, where), values
}

func (r *repo) LatestStats(ctx context.Context, suburbID, bedroom, bathroom null.Int, propertyType null.String) (entity.BuyerDemands, error) {
	resp := entity.BuyerDemands{}
	query, args := generateWhereClause(suburbID, bedroom, bathroom, propertyType)

	rows, err := r.db.Query(ctx, query, args...)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}

	defer rows.Close()
	for rows.Next() {
		v := entity.BuyerDemand{}
		err := rows.Scan(
			&v.MedianDaysToSell,
			&v.MedianSalePrice,
			&v.NumOfForSaleProperties,
			&v.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Next()")
		}
		resp = append(resp, v)
	}
	return resp, nil
}
