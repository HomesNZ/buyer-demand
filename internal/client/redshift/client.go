package redshift

import (
	"context"

	"github.com/HomesNZ/buyer-demand/internal/client/redshift/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Client interface {
	DailyBuyerDemandTablerefresh(ctx context.Context) error
}

type client struct {
	conn   *pgxpool.Pool
	config *config.Config
}

func (c client) runRefreshQuery(ctx context.Context, query string) error {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "Error while starting transaction")
	}
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "Error while executing query")
	}
	err = tx.Commit(ctx)
	return errors.Wrap(err, "Error while commiting results")
}

const dailyBuyerDemandTableRefreshQuery = ``

func (c client) DailyBuyerDemandTablerefresh(ctx context.Context) error {
	return errors.Wrap(c.runRefreshQuery(ctx, dailyBuyerDemandTableRefreshQuery), "Error while refreshing daily listing table")
}
