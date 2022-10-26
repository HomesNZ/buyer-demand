package service

import (
	"context"

	es "github.com/HomesNZ/buyer-demand/internal/client/elasticsearch"
	redshift "github.com/HomesNZ/buyer-demand/internal/client/redshift"
	"github.com/HomesNZ/gateway/errors"
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

func (s service) DailyBuyerDemandTableRefresh(ctx context.Context) error {
	resp, err := s.esClient.QueryAllListings(ctx)
	err = s.redshiftClient.DailyBuyerDemandTablerefresh(ctx)
	if err != nil {
		return errors.Wrap(err, "DailyBuyerDemandTablerefresh")
	}
	s.logger.Info("DailyBuyerDemandTablerefresh is done")
	return nil
}
