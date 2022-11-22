package service

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/api"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/HomesNZ/buyer-demand/internal/util"
	"github.com/HomesNZ/go-secret/auth"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

func (s service) BuyerDemandLatestStats(ctx context.Context, req *api.BuyerDemandLatestStatsRequest) (*api.BuyerDemandStatsResponse, error) {
	if !req.User.HasRole(auth.RoleAdmin) {
		isOwner, err := s.repos.PropertyClaim().IsClaimedByUserID(ctx, req.PropertyID, req.User.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "PropertyClaim.IsClaimedByUserID")
		}
		if !isOwner {
			return nil, util.Unauthorized(fmt.Sprintf("The property is not claimed by %d", req.User.UserID))
		}
	}

	property, err := s.esClient.ByPropertyID(ctx, req.PropertyID)
	if err != nil {
		return nil, errors.Wrap(err, "esClient.ByPropertyID")
	}

	suburbID := property.Address.SuburbID
	propertySubCategory := property.PropertySubCategory
	bedrooms := property.NumBedrooms
	bathrooms := property.NumBathrooms
	if req.NumBedrooms.Valid {
		bedrooms = req.NumBedrooms
	}
	if req.NumBathrooms.Valid {
		bathrooms = req.NumBathrooms
	}

	buyerDemand, err := s.repos.BuyerDemand().LatestStats(ctx, suburbID, bedrooms, bathrooms, propertySubCategory)
	if err != nil {
		return nil, errors.Wrap(err, "BuyerDemand.LatestStats")
	}

	stats := &api.BuyerDemandStatsResponse{
		NumBedrooms:            bedrooms,
		NumBathrooms:           bathrooms,
		SuburbID:               suburbID,
		PropertyType:           convertPropertySubCategoryToType(propertySubCategory),
		NumOfForSaleProperties: buyerDemand.NumOfForSaleProperties,
		CreatedAt:              buyerDemand.CreatedAt,
	}

	return handlerBuyerDemands(stats, buyerDemand), nil
}

func handlerBuyerDemands(stats *api.BuyerDemandStatsResponse, buyerDemand entity.BuyerDemand) *api.BuyerDemandStatsResponse {
	if buyerDemand.IsEmpty() {
		return stats
	}

	if buyerDemand.CurrentRangeMedianDaysToSell.Valid && buyerDemand.PreviousRangeMedianDaysToSell.Valid {
		medianDaysToSellTrendPercent, err := util.IncreasedPercent(buyerDemand.CurrentRangeMedianDaysToSell.ValueOrZero(), buyerDemand.PreviousRangeMedianDaysToSell.ValueOrZero(), 1)
		if err == nil {
			stats.MedianDaysToSellTrendPercent = null.FloatFrom(medianDaysToSellTrendPercent)
		}
	}

	if buyerDemand.CurrentRangeMedianSalePrice.Valid && buyerDemand.PreviousRangeMedianSalePrice.Valid {
		medianSalePriceTrendPercent, err := util.IncreasedPercent(buyerDemand.CurrentRangeMedianSalePrice.ValueOrZero(), buyerDemand.PreviousRangeMedianSalePrice.ValueOrZero(), 1)
		if err == nil {
			stats.MedianSalePriceTrendPercent = null.FloatFrom(medianSalePriceTrendPercent)
		}
	}

	if buyerDemand.CurrentRangeNumOfForSaleProperties.Valid && buyerDemand.PreviousRangeNumOfForSaleProperties.Valid {
		numOfForSalePropertiesTrendPercent, err := util.IncreasedPercent(buyerDemand.CurrentRangeNumOfForSaleProperties.ValueOrZero(), buyerDemand.PreviousRangeNumOfForSaleProperties.ValueOrZero(), 1)
		if err == nil {
			stats.NumOfForSalePropertiesTrendPercent = null.FloatFrom(numOfForSalePropertiesTrendPercent)
		}
	}

	return stats
}

func convertPropertySubCategoryToType(propertySubCategory null.String) null.String {
	if !propertySubCategory.Valid {
		return null.String{}
	}
	switch propertySubCategory.ValueOrZero() {
	case "RA":
		return null.StringFrom("apartment") // Residential Apartments
	case "RB":
		return null.StringFrom("bare_land") // Residential Bare or unimproved land
	case "RC":
		return null.StringFrom("converted_to_flats") // Residential Converted to flats
	case "RD":
		return null.StringFrom("house") // Residential Houses of a fully detached or semi- detached style
	case "RF":
		return null.StringFrom("home_units") // Residential Home units or flats
	case "RH":
		return null.StringFrom("home_and_income") // Residential Home and income
	case "RM":
		return null.StringFrom("multi_unit_bare_land") // Residential Bare land (multi unit)
	case "RN":
		return null.StringFrom("multiple_dwellings") // Residential Multiple houses on section
	case "RP":
		return null.StringFrom("parking") // Residential Parking
	case "RR":
		return null.StringFrom("rental_flat") // Residential Rental flats
	case "RV":
		return null.StringFrom("vacant_land") // Residential Vacant land
	case "LB":
		return null.StringFrom("lifestyle_bare") // Lifestyle Bare
	case "LI":
		return null.StringFrom("lifestyle_improved") // Lifestyle Improved
	case "LV":
		return null.StringFrom("lifestyle_vacant") // Lifestyle Vacant
	default:
		return null.String{}
	}
}
