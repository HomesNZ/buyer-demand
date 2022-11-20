package api

import (
	"github.com/HomesNZ/go-secret/auth"
	"gopkg.in/guregu/null.v3"
)

type BuyerDemandLatestStatsByPropertyIDRequest struct {
	PropertyID   string
	User         *auth.User
	NumBedrooms  null.Int `schema:"num_bedrooms"`
	NumBathrooms null.Int `schema:"num_bathrooms"`
}

type BuyerDemandStatsResponse struct {
	NumBedrooms                        null.Int    `json:"num_bedrooms"`
	NumBathrooms                       null.Int    `json:"num_bathrooms"`
	SuburbID                           null.Int    `json:"suburb_id"`
	PropertyType                       null.String `json:"property_type"`
	MedianDaysToSell                   null.Int    `json:"median_days_to_sell"`
	MedianSalePrice                    null.Float  `json:"median_sale_price"`
	NumOfForSaleProperties             null.Int    `json:"num_for_sale_properties"`
	MedianDaysToSellTrendPercent       null.Float  `json:"median_days_to_sell_trend_percent"`
	MedianSalePriceTrendPercent        null.Float  `json:"median_sale_price_trend_percent"`
	NumOfForSalePropertiesTrendPercent null.Float  `json:"num_for_sale_properties_trend_percent"`
	CreatedAt                          null.Time   `json:"created_at"`
}
