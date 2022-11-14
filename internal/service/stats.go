package service

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/api"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/HomesNZ/buyer-demand/internal/util"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
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

	suburbID := property.Address.SuburbID
	bedrooms := property.NumBedrooms
	bathrooms := property.NumBathrooms
	propertySubCategory := property.PropertySubCategory

	buyerDemands, err := s.repos.BuyerDemand().LatestStats(ctx, suburbID, bedrooms, bathrooms, propertySubCategory)
	if err != nil {
		return nil, errors.Wrap(err, "BuyerDemand.LatestStats")
	}

	stats := &api.BuyerDemandStatsResponse{
		NumBedrooms:  bedrooms,
		NumBathrooms: bathrooms,
		SuburbID:     suburbID,
		PropertyType: convertPropertySubCategoryToType(propertySubCategory),
	}

	return handlerBuyerDemands(stats, buyerDemands), nil
}

func handlerBuyerDemands(stats *api.BuyerDemandStatsResponse, buyerDemands entity.BuyerDemands) *api.BuyerDemandStatsResponse {
	if len(buyerDemands) == 0 {
		return stats
	}

	if len(buyerDemands) >= 1 {
		stats.MedianDaysToSell = buyerDemands[0].MedianDaysToSell
		stats.MedianSalePrice = buyerDemands[0].MedianSalePrice
		stats.NumOfForSaleProperties = buyerDemands[0].NumOfForSaleProperties
		stats.CreatedAt = buyerDemands[0].CreatedAt
	}

	if len(buyerDemands) >= 2 {
		if buyerDemands[0].MedianDaysToSell.Valid && buyerDemands[1].MedianDaysToSell.Valid {
			medianDaysToSellIncreasedPercent, err := util.IncreasedPercent(buyerDemands[0].MedianDaysToSell.ValueOrZero(), buyerDemands[1].MedianDaysToSell.ValueOrZero())
			if err == nil {
				stats.MedianDaysToSellIncreasedPercent = null.FloatFrom(medianDaysToSellIncreasedPercent)
			}
		}

		if buyerDemands[0].MedianSalePrice.Valid && buyerDemands[1].MedianSalePrice.Valid {
			medianSalePriceIncreasedPercent, err := util.IncreasedPercent(buyerDemands[0].MedianSalePrice.ValueOrZero(), buyerDemands[1].MedianSalePrice.ValueOrZero())
			if err == nil {
				stats.MedianSalePriceIncreasedPercent = null.FloatFrom(medianSalePriceIncreasedPercent)
			}
		}

		if buyerDemands[0].NumOfForSaleProperties.Valid && buyerDemands[1].NumOfForSaleProperties.Valid {
			numOfForSalePropertiesIncreasedPercent, err := util.IncreasedPercent(buyerDemands[0].NumOfForSaleProperties.ValueOrZero(), buyerDemands[1].NumOfForSaleProperties.ValueOrZero())
			if err == nil {
				stats.NumOfForSalePropertiesIncreasedPercent = null.FloatFrom(numOfForSalePropertiesIncreasedPercent)
			}
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
