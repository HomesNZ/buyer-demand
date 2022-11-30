package service

import (
	"encoding/json"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v3"
	"sort"
	"strings"
	"time"
)

var _ = Describe("BuyerDemand", func() {

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 1, 2, 3, 4, now.Location())
	currentRange := today.AddDate(0, 0, -90)
	previousRange := currentRange.AddDate(0, 0, -180)
	outOfRange := previousRange.AddDate(0, 0, -180)

	address1020 := entity.Address{
		SuburbID: null.IntFrom(1020),
	}
	addressNull := entity.Address{
		SuburbID: null.Int{},
	}
	mapItems := entity.MapItemESs{
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(600000),
			ListingId:           null.StringFrom("listing-test-1"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(900000),
			ListingId:           null.StringFrom("listing-test-2"),
			LatestListingDate:   null.TimeFrom(now.AddDate(0, 0, -35)),
			LatestSoldDate:      null.Time{},
		},
		entity.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(650000),
			ListingId:           null.StringFrom("listing-test-3"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.Int{},
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(850000),
			ListingId:           null.StringFrom("listing-test-4"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             addressNull,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(725000),
			ListingId:           null.StringFrom("listing-test-5"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.TimeFrom(now.AddDate(-1, 0, -3)),
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.String{},
			Price:               null.FloatFrom(750000),
			ListingId:           null.StringFrom("listing-test-6"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.TimeFrom(now.AddDate(0, -2, -1)),
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RP"),
			Price:               null.Float{},
			ListingId:           null.StringFrom("listing-test-7"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.TimeFrom(now.AddDate(0, 0, -10)),
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RP"),
			Price:               null.FloatFrom(800000),
			ListingId:           null.StringFrom("listing-test-8"),
			LatestListingDate:   null.Time{},
			LatestSoldDate:      null.Time{},
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(753000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(currentRange),
			LatestSoldDate:      null.TimeFrom(today),
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(853000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(previousRange),
			LatestSoldDate:      null.TimeFrom(currentRange),
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(830000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(outOfRange),
			LatestSoldDate:      null.TimeFrom(previousRange),
		},
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(630000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(outOfRange),
			LatestSoldDate:      null.TimeFrom(outOfRange),
		},
		entity.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(800000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(previousRange),
			LatestSoldDate:      null.TimeFrom(outOfRange),
		},
		entity.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(830000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(previousRange),
			LatestSoldDate:      null.TimeFrom(currentRange),
		},
		entity.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(760000),
			ListingId:           null.String{},
			PropertyState:       null.Int{},
			LatestListingDate:   null.TimeFrom(outOfRange),
			LatestSoldDate:      null.TimeFrom(previousRange),
		},
	}

	bds := entity.BuyerDemands{
		entity.BuyerDemand{
			NumBedrooms:                         null.Int{},
			NumBathrooms:                        null.IntFrom(3),
			SuburbID:                            null.IntFrom(1020),
			PropertyType:                        null.StringFrom("RR"),
			CurrentRangeMedianDaysToSell:        null.IntFrom(180),
			PreviousRangeMedianDaysToSell:       null.Int{},
			CurrentRangeMedianSalePrice:         null.FloatFrom(830000),
			PreviousRangeMedianSalePrice:        null.Float{},
			NumOfForSaleProperties:              null.IntFrom(1),
			CurrentRangeNumOfForSaleProperties:  null.IntFrom(1),
			PreviousRangeNumOfForSaleProperties: null.IntFrom(0),
		},
		entity.BuyerDemand{
			NumBedrooms:                         null.IntFrom(4),
			NumBathrooms:                        null.Int{},
			SuburbID:                            null.IntFrom(1020),
			PropertyType:                        null.StringFrom("RR"),
			CurrentRangeMedianDaysToSell:        null.Int{},
			PreviousRangeMedianDaysToSell:       null.Int{},
			CurrentRangeMedianSalePrice:         null.Float{},
			PreviousRangeMedianSalePrice:        null.Float{},
			NumOfForSaleProperties:              null.IntFrom(1),
			CurrentRangeNumOfForSaleProperties:  null.IntFrom(1),
			PreviousRangeNumOfForSaleProperties: null.IntFrom(0),
		},
		entity.BuyerDemand{
			NumBedrooms:                         null.IntFrom(4),
			NumBathrooms:                        null.IntFrom(3),
			SuburbID:                            null.Int{},
			PropertyType:                        null.StringFrom("RR"),
			CurrentRangeMedianDaysToSell:        null.Int{},
			PreviousRangeMedianDaysToSell:       null.Int{},
			CurrentRangeMedianSalePrice:         null.Float{},
			PreviousRangeMedianSalePrice:        null.Float{},
			NumOfForSaleProperties:              null.IntFrom(1),
			CurrentRangeNumOfForSaleProperties:  null.IntFrom(1),
			PreviousRangeNumOfForSaleProperties: null.IntFrom(0),
		},
		entity.BuyerDemand{
			NumBedrooms:                         null.IntFrom(4),
			NumBathrooms:                        null.IntFrom(3),
			SuburbID:                            null.IntFrom(1020),
			PropertyType:                        null.String{},
			CurrentRangeMedianDaysToSell:        null.Int{},
			PreviousRangeMedianDaysToSell:       null.Int{},
			CurrentRangeMedianSalePrice:         null.Float{},
			PreviousRangeMedianSalePrice:        null.Float{},
			NumOfForSaleProperties:              null.IntFrom(1),
			CurrentRangeNumOfForSaleProperties:  null.IntFrom(1),
			PreviousRangeNumOfForSaleProperties: null.IntFrom(0),
		},
		entity.BuyerDemand{
			NumBedrooms:                         null.IntFrom(4),
			NumBathrooms:                        null.IntFrom(3),
			SuburbID:                            null.IntFrom(1020),
			PropertyType:                        null.StringFrom("RR"),
			CurrentRangeMedianDaysToSell:        null.IntFrom(135),
			PreviousRangeMedianDaysToSell:       null.IntFrom(180),
			CurrentRangeMedianSalePrice:         null.FloatFrom(803000),
			PreviousRangeMedianSalePrice:        null.FloatFrom(830000),
			NumOfForSaleProperties:              null.IntFrom(2),
			CurrentRangeNumOfForSaleProperties:  null.IntFrom(1),
			PreviousRangeNumOfForSaleProperties: null.IntFrom(1),
		},
		entity.BuyerDemand{
			NumBedrooms:                         null.IntFrom(4),
			NumBathrooms:                        null.IntFrom(3),
			SuburbID:                            null.IntFrom(1020),
			PropertyType:                        null.StringFrom("RP"),
			CurrentRangeMedianDaysToSell:        null.Int{},
			PreviousRangeMedianDaysToSell:       null.Int{},
			CurrentRangeMedianSalePrice:         null.Float{},
			PreviousRangeMedianSalePrice:        null.Float{},
			NumOfForSaleProperties:              null.IntFrom(2),
			CurrentRangeNumOfForSaleProperties:  null.IntFrom(1),
			PreviousRangeNumOfForSaleProperties: null.IntFrom(0),
		},
	}

	Describe("DailyBuyerDemandTableRefresh", func() {

		It("generateBuyerDemand", func() {
			bdsActual := mapItems.GenerateBuyerDemands()

			sort.Slice(bdsActual, func(i, j int) bool {
				ii, _ := json.Marshal(bdsActual[i])
				jj, _ := json.Marshal(bdsActual[j])
				return strings.Compare(string(ii), string(jj)) > 0
			})

			bdsActualJson, _ := json.Marshal(bdsActual)
			bdsJson, _ := json.Marshal(bds)

			Expect(string(bdsJson)).To(Equal(string(bdsActualJson)))
		})
	})
})
