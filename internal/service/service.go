package service

import (
	"context"

	es "github.com/HomesNZ/buyer-demand/internal/client/elasticsearch"
	"github.com/HomesNZ/buyer-demand/internal/client/redshift"
	"github.com/sirupsen/logrus"
)

func New(log *logrus.Entry, redshift redshift.Client, esClient es.Client) (Service, error) {
	s := &service{
		redshiftClient: redshift,
		esClient:       esClient,
		logger:         log,
	}

	return s, nil
}

type Service interface {
	DailyBuyerDemandTableRefresh(ctx context.Context) error
}

type service struct {
	redshiftClient redshift.Client
	esClient       es.Client
	logger         *logrus.Entry
}
