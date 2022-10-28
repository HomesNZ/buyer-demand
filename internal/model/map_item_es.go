package model

import (
	"github.com/HomesNZ/buyer-demand/internal/util"
	"gopkg.in/guregu/null.v3"
	"math"
	"strconv"
	"strings"
	"time"
)

type MapItemES struct {
	NumBedrooms           null.Int    `json:"num_bedrooms"`
	NumBathrooms          null.Int    `json:"num_bathrooms"`
	Suburb                null.String `json:"nested_address.suburb"`
	CategoryCode          null.String `json:"category_code"`
	Price                 null.Float  `json:"price"`
	ListingId             null.String `json:"listing_id"`
	ListingRecentSaleDate null.Time   `json:"listing_recent_sale_date"`
}

type MapItemESs []MapItemES

func (i *MapItemES) getKey() buyerDemandKey {
	return buyerDemandKey(strings.Join(
		[]string{
			strconv.FormatInt(i.NumBedrooms.ValueOrZero(), 10),
			strconv.FormatInt(i.NumBathrooms.ValueOrZero(), 10),
			i.Suburb.ValueOrZero(),
			i.CategoryCode.ValueOrZero()[0:2],
		}, BuyerDemandKeySeparator))
}

func (items MapItemESs) GenerateBuyerDemands() BuyerDemands {
	currentListingMap, daysToSellMap := items.prepareData()

	var result BuyerDemands
	for key, currentListings := range currentListingMap {
		bd := key.generateBuyerDemandFromKey()

		var medianDaysToSell null.Int
		daysToSell, ok := daysToSellMap[key]
		if ok {
			medianDaysToSell = calculateMedianDaysToSell(daysToSell)
		}
		bd.MedianDaysToSell = medianDaysToSell
		bd.MedianSalePrice = currentListings.calculateMedianSalePrice()
		bd.NumOfForSaleProperties = currentListings.calculateNumOfForSaleProperties()

		result = append(result, bd)
	}

	return result
}

func (items MapItemESs) prepareData() (map[buyerDemandKey]MapItemESs, map[buyerDemandKey][]int64) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastNinetyDays := today.AddDate(0, 0, 90)

	var currentListingMap map[buyerDemandKey]MapItemESs
	var daysToSellMap map[buyerDemandKey][]int64

	for _, item := range items {
		key := item.getKey()
		if item.ListingId.Valid {
			currentListings, ok := currentListingMap[key]
			if !ok {
				currentListings = MapItemESs{}
			}

			currentListings = append(currentListings, item)
			currentListingMap[key] = currentListings

			continue
		}

		if item.ListingRecentSaleDate.Valid && item.ListingRecentSaleDate.Time.After(lastNinetyDays) {
			daysToSell, ok := daysToSellMap[key]
			if !ok {
				daysToSell = []int64{}
			}

			t := item.ListingRecentSaleDate.ValueOrZero()
			d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			days := int64(math.Round(today.Sub(d).Hours() / 24))

			daysToSell = append(daysToSell, days)
			daysToSellMap[key] = daysToSell

			continue
		}
	}

	return currentListingMap, daysToSellMap
}

func calculateMedianDaysToSell(daysToSell []int64) null.Int {
	return null.IntFrom(int64(util.Median(daysToSell)))
}

func (items MapItemESs) calculateMedianSalePrice() float64 {
	var priceArray []float64
	for _, item := range items {
		priceArray = append(priceArray, item.Price.ValueOrZero())
	}

	return util.Median(priceArray)
}

func (items MapItemESs) calculateNumOfForSaleProperties() int {
	return len(items)
}
