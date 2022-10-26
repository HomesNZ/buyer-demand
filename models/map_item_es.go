package mapItemES

import (
	"gopkg.in/guregu/null.v3"
)

type MapItemES struct {
	NumBedrooms  null.Int    `json:"num_bedrooms"`
	NumBathrooms null.Int    `json:"num_bathrooms"`
	Suburb       null.String `json:"nested_address.suburb"`
	CategoryCode null.String `json:"category_code"`
}

type MapItemESs []*MapItemES
