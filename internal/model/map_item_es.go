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
	NumBedrooms         null.Int    `json:"num_bedrooms"`
	NumBathrooms        null.Int    `json:"num_bathrooms"`
	Suburb              null.String `json:"nested_address.suburb"`
	PropertySubCategory null.String `json:"property_sub_category"`
	Price               null.Float  `json:"price"`
	ListingId           null.String `json:"listing_id"`
	PropertyState       null.Int    `json:"property_state"`
	LatestListingDate   null.Time   `json:"latest_listing_date"`
	LatestSoldDate      null.Time   `json:"latest_sold_date"`
}

type MapItemESs []MapItemES

func (i *MapItemES) getKey() buyerDemandKey {
	return buyerDemandKey(strings.Join(
		[]string{
			strconv.FormatInt(i.NumBedrooms.ValueOrZero(), 10),
			strconv.FormatInt(i.NumBathrooms.ValueOrZero(), 10),
			i.Suburb.ValueOrZero(),
			i.PropertySubCategory.ValueOrZero(),
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
	lastNinetyDays := today.AddDate(0, 0, -90)

	currentListingMap := map[buyerDemandKey]MapItemESs{}
	daysToSellMap := map[buyerDemandKey][]int64{}

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

		if item.PropertyState.Valid && item.PropertyState.ValueOrZero() == 2 &&
			item.LatestListingDate.Valid && item.LatestSoldDate.Valid &&
			item.LatestSoldDate.ValueOrZero().After(item.LatestListingDate.ValueOrZero()) &&
			item.LatestListingDate.Time.After(lastNinetyDays) {

			daysToSell, ok := daysToSellMap[key]
			if !ok {
				daysToSell = []int64{}
			}

			listingSoldTime := item.LatestSoldDate.ValueOrZero()
			listingSoldDate := time.Date(listingSoldTime.Year(), listingSoldTime.Month(), listingSoldTime.Day(), 0, 0, 0, 0, listingSoldTime.Location())

			listingTime := item.LatestListingDate.ValueOrZero()
			listingDate := time.Date(listingTime.Year(), listingTime.Month(), listingTime.Day(), 0, 0, 0, 0, listingTime.Location())
			days := int64(math.Round(listingSoldDate.Sub(listingDate).Hours() / 24))

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
