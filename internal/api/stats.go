package api

import (
	"gopkg.in/guregu/null.v3"
)

type BuyerDemandLatestStatsByPropertyIDRequest struct {
	PropertyID string
	UserID     int
}

type BuyerDemandStatsResponse struct {
	MedianDaysToSell       null.Int   `json:"median_days_to_sell"`
	MedianSalePrice        null.Float `json:"median_sale_price"`
	NumOfForSaleProperties null.Int   `json:"num_for_sale_properties"`
	CreatedAt              null.Time  `json:"created_at"`
}
