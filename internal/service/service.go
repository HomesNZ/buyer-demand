package service

import (
	"context"
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
	DailyBuyerDemandTableRefresh(ctx context.Context) error
}

type service struct {
	repos    repository.Repositories
	esClient es.Client
	logger   *logrus.Entry
}
