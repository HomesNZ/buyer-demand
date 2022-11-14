package entity

import (
	"gopkg.in/guregu/null.v3"
	"strconv"
	"strings"
)

const BuyerDemandKeySeparator = "-"

type BuyerDemand struct {
	NumBedrooms            null.Int    `json:"num_bedrooms"`
	NumBathrooms           null.Int    `json:"num_bathrooms"`
	SuburbID               null.Int    `json:"suburb_id"`
	PropertyType           null.String `json:"property_type"`
	MedianDaysToSell       null.Int    `json:"median_days_to_sell"`
	MedianSalePrice        null.Float  `json:"median_sale_price"`
	NumOfForSaleProperties null.Int    `json:"num_for_sale_properties"`
	CreatedAt              null.Time   `json:"created_at"`
}

type BuyerDemands []BuyerDemand

type buyerDemandKey string

func (key buyerDemandKey) generateBuyerDemandFromKey() BuyerDemand {
	keyArray := strings.Split(string(key), BuyerDemandKeySeparator)

	part1, err := strconv.ParseInt(keyArray[0], 0, 64)
	if err != nil {
		part1 = 0
	}
	var numBedroom null.Int
	if part1 != 0 {
		numBedroom = null.IntFrom(part1)
	}

	part2, err := strconv.ParseInt(keyArray[1], 0, 64)
	if err != nil {
		part2 = 0
	}
	var numBathroom null.Int
	if part2 != 0 {
		numBathroom = null.IntFrom(part2)
	}

	part3, err := strconv.ParseInt(keyArray[2], 0, 64)
	if err != nil {
		part3 = 0
	}
	var suburbID null.Int
	if part3 != 0 {
		suburbID = null.IntFrom(part3)
	}

	var propertyType null.String
	if keyArray[3] != "" {
		propertyType = null.StringFrom(keyArray[3])
	}

	return BuyerDemand{
		NumBedrooms:  numBedroom,
		NumBathrooms: numBathroom,
		SuburbID:     suburbID,
		PropertyType: propertyType,
	}
}
