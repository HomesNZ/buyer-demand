package model

import (
	"gopkg.in/guregu/null.v3"
	"strconv"
	"strings"
)

const BuyerDemandKeySeparator = "-"

type BuyerDemand struct {
	NumBedrooms            null.Int    `json:"num_bedrooms"`
	NumBathrooms           null.Int    `json:"num_bathrooms"`
	Suburb                 null.String `json:"nested_address.suburb"`
	PropertyType           null.String `json:"property_type"`
	MedianDaysToSell       null.Int    `json:"median_days_to_sell"`
	MedianSalePrice        float64     `json:"median_sale_price"`
	NumOfForSaleProperties int         `json:"num_of_for_sale_properties"`
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

	var suburb null.String
	if keyArray[2] != "" {
		suburb = null.StringFrom(keyArray[2])
	}

	var propertyType null.String
	if keyArray[3] != "" {
		propertyType = null.StringFrom(keyArray[3])
	}

	return BuyerDemand{
		NumBedrooms:  numBedroom,
		NumBathrooms: numBathroom,
		Suburb:       suburb,
		PropertyType: propertyType,
	}
}
