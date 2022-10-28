package service

import (
	"context"
	"github.com/pkg/errors"
)

func (s service) DailyBuyerDemandTableRefresh(ctx context.Context) error {
	mapItems, err := s.esClient.QueryAll(ctx)
	if err != nil {
		return errors.Wrap(err, "QueryAll")
	}
	bds := mapItems.GenerateBuyerDemands()
	err = s.redshiftClient.DailyBuyerDemandTableRefresh(ctx, bds)
	if err != nil {
		return errors.Wrap(err, "DailyBuyerDemandTableRefresh")
	}
	s.logger.Info("DailyBuyerDemandTableRefresh is done")
	return nil
}
