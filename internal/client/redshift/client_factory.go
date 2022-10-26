package redshift

import (
	"context"

	"github.com/HomesNZ/buyer-demand/internal/client/redshift/config"
	"github.com/HomesNZ/go-common/dbclient/v4"
	dbConfig "github.com/HomesNZ/go-common/dbclient/v4/config"
	"github.com/pkg/errors"
)

func NewFromEnv(ctx context.Context) (Client, error) {
	cfg, err := config.NewFromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to extract environment variables")
	}
	db, err := dbclient.New(ctx, &dbConfig.Config{
		Port:        cfg.Port,
		Password:    cfg.Password,
		User:        cfg.User,
		Host:        cfg.Host,
		Name:        cfg.Name,
		ServiceName: cfg.ServiceName,
		MaxConns:    cfg.MaxConns,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to Redshift")
	}

	return &client{conn: db, config: cfg}, nil
}
