package service

import (
	"context"
	"github.com/HomesNZ/buyer-demand/internal/api"
	"github.com/HomesNZ/buyer-demand/internal/repository"

	es "github.com/HomesNZ/buyer-demand/internal/client/elasticsearch"
	"github.com/sirupsen/logrus"
)

func New(log *logrus.Entry, repos repository.Repositories, esClient es.Client) (Service, error) {
	s := &service{
		repos:    repos,
		esClient: esClient,
		logger:   log,
	}

	return s, nil
}

type Service interface {
	Health() error
	DailyBuyerDemandTableRefresh(ctx context.Context) error
	BuyerDemandLatestStats(ctx context.Context, req *api.BuyerDemandStatsRequest) (*api.BuyerDemandStatsResponse, error)
}

type service struct {
	repos    repository.Repositories
	esClient es.Client
	logger   *logrus.Entry
}

func (s service) Health() error {
	return nil
}
