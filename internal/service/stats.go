package service

import (
	"context"
	"github.com/HomesNZ/buyer-demand/internal/api"
	"github.com/pkg/errors"
)

func (s service) BuyerDemandLatestStats(ctx context.Context, req *api.BuyerDemandStatsRequest) (*api.BuyerDemandStatsResponse, error) {
	if req.SuburbID.ValueOrZero() == 0 {
		return nil, errors.New("Invalid SuburbID")
	}

	stats, err := s.repos.BuyerDemand().LatestStats(ctx, req.SuburbID, req.NumBedrooms, req.NumBathrooms, req.PropertyType)
	if err != nil {
		return nil, errors.Wrap(err, "BuyerDemand.LatestStats")
	}

	return stats, nil
}
