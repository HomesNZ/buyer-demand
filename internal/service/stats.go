package service

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/api"
	"github.com/HomesNZ/buyer-demand/internal/util"
	"github.com/pkg/errors"
)

func (s service) BuyerDemandLatestStatsByPropertyID(ctx context.Context, req *api.BuyerDemandLatestStatsByPropertyIDRequest) (*api.BuyerDemandStatsResponse, error) {
	isOwner, err := s.repos.PropertyClaim().IsClaimedByUserID(ctx, req.PropertyID, req.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "PropertyClaim.IsClaimedByUserID")
	}
	if !isOwner {
		return nil, util.Unauthorized(fmt.Sprintf("The property is not claimed by %d", req.UserID))
	}

	property, err := s.esClient.ByPropertyID(ctx, req.PropertyID)
	if err != nil {
		return nil, errors.Wrap(err, "esClient.ByPropertyID")
	}

	stats, err := s.repos.BuyerDemand().LatestStats(ctx, property.Address.SuburbID, property.NumBedrooms, property.NumBathrooms, property.PropertySubCategory)
	if err != nil {
		return nil, errors.Wrap(err, "BuyerDemand.LatestStats")
	}

	return stats, nil
}
