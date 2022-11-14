package api

import (
	"gopkg.in/guregu/null.v3"
)

type BuyerDemandLatestStatsByPropertyIDRequest struct {
	PropertyID string
	UserID     int
}

type BuyerDemandStatsResponse struct {
	NumBedrooms                            null.Int    `json:"num_bedrooms"`
	NumBathrooms                           null.Int    `json:"num_bathrooms"`
	SuburbID                               null.Int    `json:"suburb_id"`
	PropertyType                           null.String `json:"property_type"`
	MedianDaysToSell                       null.Int    `json:"median_days_to_sell"`
	MedianSalePrice                        null.Float  `json:"median_sale_price"`
	NumOfForSaleProperties                 null.Int    `json:"num_for_sale_properties"`
	MedianDaysToSellIncreasedPercent       null.Float  `json:"median_days_to_sell_increased_percent"`
	MedianSalePriceIncreasedPercent        null.Float  `json:"median_sale_price_increased_percent"`
	NumOfForSalePropertiesIncreasedPercent null.Float  `json:"num_for_sale_properties_increased_percent"`
	CreatedAt                              null.Time   `json:"created_at"`
}
