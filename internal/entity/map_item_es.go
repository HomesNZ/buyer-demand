package entity

import (
	"github.com/HomesNZ/buyer-demand/internal/util"
	"gopkg.in/guregu/null.v3"
	"math"
	"strconv"
	"strings"
	"time"
)

type MapItemES struct {
	NumBedrooms         null.Int    `json:"num_bedrooms"`
	NumBathrooms        null.Int    `json:"num_bathrooms"`
	Address             Address     `json:"nested_address"`
	PropertySubCategory null.String `json:"property_sub_category"`
	Price               null.Float  `json:"price"`
	ListingId           null.String `json:"listing_id"`
	PropertyState       null.Int    `json:"property_state"`
	LatestListingDate   null.Time   `json:"latest_listing_date"`
	LatestSoldDate      null.Time   `json:"latest_sold_date"`
}

type Address struct {
	SuburbID null.Int `json:"suburb_id"`
}

type MapItemESs []MapItemES

type BuyerDemandES struct {
	numberForSaleProperties              int
	currentRangeNumberForSaleProperties  int
	previousRangeNumberForSaleProperties int
	currentRangeDaysToSellList           []int64
	previousRangeDaysToSellList          []int64
	currentRangeSalePriceList            []float64
	previousRangeSalePriceList           []float64
}

func (bd *BuyerDemandES) appendNumberForSaleProperties(i *MapItemES, currentRangeStartDate time.Time, previousRangeStartDate time.Time) *BuyerDemandES {
	if i.isCurrentListing() {
		bd.numberForSaleProperties++
	}

	listingDate := util.ToUtcDate(i.LatestListingDate.ValueOrZero())
	if listingDate.After(currentRangeStartDate) {
		bd.currentRangeNumberForSaleProperties++
	} else if listingDate.After(previousRangeStartDate) {
		bd.previousRangeNumberForSaleProperties++
	}

	return bd
}

func (bd *BuyerDemandES) appendDaysToSell(i *MapItemES, currentRangeStartDate time.Time, previousRangeStartDate time.Time) *BuyerDemandES {
	if i.LatestListingDate.IsZero() || i.LatestSoldDate.IsZero() {
		return bd
	}

	listingDate := util.ToUtcDate(i.LatestListingDate.ValueOrZero())
	listingSoldDate := util.ToUtcDate(i.LatestSoldDate.ValueOrZero())

	if listingDate.After(listingSoldDate) || previousRangeStartDate.After(listingSoldDate) {
		return bd
	}

	days := int64(math.Round(listingSoldDate.Sub(listingDate).Hours() / 24))

	if listingSoldDate.After(currentRangeStartDate) {
		// current time frame range
		daysToSellList := bd.currentRangeDaysToSellList
		daysToSellList = append(daysToSellList, days)
		bd.currentRangeDaysToSellList = daysToSellList
	} else {
		// previous time frame range
		daysToSellList := bd.previousRangeDaysToSellList
		daysToSellList = append(daysToSellList, days)
		bd.previousRangeDaysToSellList = daysToSellList
	}

	return bd
}

func (bd *BuyerDemandES) appendSalePrice(i *MapItemES, currentRangeStartDate time.Time, previousRangeStartDate time.Time) *BuyerDemandES {
	if i.Price.IsZero() || i.LatestSoldDate.IsZero() {
		return bd
	}

	soldDate := util.ToUtcDate(i.LatestSoldDate.ValueOrZero())
	if soldDate.After(currentRangeStartDate) {
		salePriceList := bd.currentRangeSalePriceList
		salePriceList = append(salePriceList, i.Price.ValueOrZero())
		bd.currentRangeSalePriceList = salePriceList
	} else if soldDate.After(previousRangeStartDate) {
		salePriceList := bd.previousRangeSalePriceList
		salePriceList = append(salePriceList, i.Price.ValueOrZero())
		bd.previousRangeSalePriceList = salePriceList
	}

	return bd
}

func (i *MapItemES) getKey() buyerDemandKey {
	return buyerDemandKey(strings.Join(
		[]string{
			strconv.FormatInt(i.NumBedrooms.ValueOrZero(), 10),
			strconv.FormatInt(i.NumBathrooms.ValueOrZero(), 10),
			strconv.FormatInt(i.Address.SuburbID.ValueOrZero(), 10),
			i.PropertySubCategory.ValueOrZero(),
		}, BuyerDemandKeySeparator))
}

func (items MapItemESs) GenerateBuyerDemands() BuyerDemands {
	var result BuyerDemands
	if len(items) == 0 {
		return nil
	}

	buyerDemandES := items.prepareData()
	for key, bdES := range buyerDemandES {
		bd := key.generateBuyerDemandFromKey()
		bd.CurrentRangeMedianDaysToSell = calculateMedianDaysToSell(bdES.currentRangeDaysToSellList)
		bd.PreviousRangeMedianDaysToSell = calculateMedianDaysToSell(bdES.previousRangeDaysToSellList)
		bd.CurrentRangeMedianSalePrice = calculateMedianSalePrice(bdES.currentRangeSalePriceList)
		bd.PreviousRangeMedianSalePrice = calculateMedianSalePrice(bdES.previousRangeSalePriceList)
		bd.NumOfForSaleProperties = null.IntFrom(int64(bdES.numberForSaleProperties))
		bd.CurrentRangeNumOfForSaleProperties = null.IntFrom(int64(bdES.currentRangeNumberForSaleProperties))
		bd.PreviousRangeNumOfForSaleProperties = null.IntFrom(int64(bdES.previousRangeNumberForSaleProperties))

		if !bd.IsEmpty() {
			result = append(result, bd)
		}
	}

	return result
}

func (items MapItemESs) prepareData() map[buyerDemandKey]*BuyerDemandES {
	now := time.Now()
	today := util.ToUtcDate(now)
	currentRangeStartDate := today.AddDate(0, 0, -180)
	previousRangeStartDate := currentRangeStartDate.AddDate(0, 0, -180)
	currentRangeStartDateForNumberForSaleProperties := today.AddDate(0, 0, -30)
	previousRangeStartDateForNumberForSaleProperties := currentRangeStartDateForNumberForSaleProperties.AddDate(0, 0, -30)

	buyerDemandESMap := map[buyerDemandKey]*BuyerDemandES{}

	for _, item := range items {
		key := item.getKey()
		if item.isRecentListing(previousRangeStartDateForNumberForSaleProperties) {
			buyerDemandES := buyerDemandESMap[key]
			if buyerDemandES == nil {
				buyerDemandES = &BuyerDemandES{}
			}
			buyerDemandESMap[key] = buyerDemandES.appendNumberForSaleProperties(&item, currentRangeStartDateForNumberForSaleProperties, previousRangeStartDateForNumberForSaleProperties)
		}

		if item.isSold() {
			buyerDemandES := buyerDemandESMap[key]
			if buyerDemandES == nil {
				buyerDemandES = &BuyerDemandES{}
			}

			buyerDemandESMap[key] = buyerDemandES.appendDaysToSell(&item, currentRangeStartDate, previousRangeStartDate)
			buyerDemandESMap[key] = buyerDemandES.appendSalePrice(&item, currentRangeStartDate, previousRangeStartDate)

			continue
		}
	}

	return buyerDemandESMap
}

func (i *MapItemES) isCurrentListing() bool {
	return i.ListingId.Valid
}

func (i *MapItemES) isRecentListing(startDate time.Time) bool {
	return (i.LatestListingDate.Valid && i.LatestListingDate.ValueOrZero().After(startDate)) || i.ListingId.Valid
}

func (i *MapItemES) isSold() bool {
	return i.PropertyState.Valid && i.PropertyState.ValueOrZero() == 2
}

func calculateMedianDaysToSell(daysToSell []int64) null.Int {
	if len(daysToSell) == 0 {
		return null.Int{}
	}

	return null.IntFrom(int64(util.Median(daysToSell)))
}

func calculateMedianSalePrice(salePrices []float64) null.Float {
	if len(salePrices) == 0 {
		return null.Float{}
	}
	return null.FloatFrom(util.Median(salePrices))
}
