package redshift

import (
	"context"
	"github.com/HomesNZ/buyer-demand/internal/model"
	"time"

	"github.com/HomesNZ/buyer-demand/internal/client/redshift/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Client interface {
	DailyBuyerDemandTableRefresh(ctx context.Context, bds model.BuyerDemands) error
}

type client struct {
	conn   *pgxpool.Pool
	config *config.Config
}

func (c client) runQueries(ctx context.Context, query string, argumentsArray [][]interface{}, deleteQuery string) error {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return errors.Wrap(err, "Error while rolling back")
		}
		return errors.Wrap(err, "Error while starting transaction")
	}

	// Delete current day data
	_, err = tx.Exec(ctx, deleteQuery)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return errors.Wrap(err, "Error while rolling back")
		}
		return errors.Wrap(err, "Error while executing query")
	}

	// Insert current day data
	for _, arguments := range argumentsArray {
		_, err = tx.Exec(ctx, query, arguments...)
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return errors.Wrap(err, "Error while rolling back")
			}
			return errors.Wrap(err, "Error while executing query")
		}
	}
	err = tx.Commit(ctx)
	return errors.Wrap(err, "Error while committing results")
}

const dailyBuyerDemandTableDeleteQuery = `
	DELETE FROM buyer_demand
	WHERE created_at >= CURRENT_DATE;
`

const dailyBuyerDemandTablePopulateQuery = `
	INSERT INTO buyer_demand (
		num_bedrooms,
	    num_bathrooms,
	    suburb,
	    property_type,
	    median_days_to_sell,
	    median_sale_price,
	    num_of_for_sale_properties,
	    created_at
	) VALUES (
	    $1, $2, $3, $4, $5, $6, $7, $8
	);
`

func (c client) DailyBuyerDemandTableRefresh(ctx context.Context, bds model.BuyerDemands) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var argumentsArray [][]interface{}
	for _, bd := range bds {
		arguments := append(make([]interface{}, 0),
			bd.NumBedrooms,
			bd.NumBathrooms,
			bd.Suburb,
			bd.PropertyType,
			bd.MedianDaysToSell,
			bd.MedianSalePrice,
			bd.NumOfForSaleProperties,
			today,
		)
		argumentsArray = append(argumentsArray, arguments)
	}

	err := c.runQueries(ctx, dailyBuyerDemandTablePopulateQuery, argumentsArray, dailyBuyerDemandTableDeleteQuery)
	if err != nil {
		return errors.Wrap(err, "Error while refreshing daily buyer demand table")
	}

	return nil
}
