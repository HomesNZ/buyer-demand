package repository

import (
	"context"
	"github.com/HomesNZ/buyer-demand/internal/repository/address"
	buyerDemand "github.com/HomesNZ/buyer-demand/internal/repository/buyer_demand"
	"github.com/HomesNZ/go-common/dbclient/v4"
	"github.com/pkg/errors"
)

func New(ctx context.Context) (Repositories, error) {
	conn, err := dbclient.NewFromEnv(ctx)
	if err != nil || conn == nil {
		return nil, errors.Wrapf(err, "failed to connect to database")
	}

	addressRepo, err := address.New(conn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create address repository")
	}

	buyerDemandRepo, err := buyerDemand.New(conn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create buyer demand repository")
	}

	return &repositories{
		address:     addressRepo,
		buyerDemand: buyerDemandRepo,
	}, nil
}
