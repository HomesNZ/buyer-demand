package address

import (
	"context"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repo interface {
	Populate(ctx context.Context, buyerDemands entity.BuyerDemands) error
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

func (r *repo) Populate(ctx context.Context, buyerDemands entity.BuyerDemands) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return errors.Wrap(err, "tx.Rollback")
		}
		return errors.Wrap(err, "db.Begin")
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
