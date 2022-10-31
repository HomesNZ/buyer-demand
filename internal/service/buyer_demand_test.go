package service

import (
	"context"
	"encoding/json"
	mockES "github.com/HomesNZ/buyer-demand/internal/client/elasticsearch/mock"
	mockRD "github.com/HomesNZ/buyer-demand/internal/client/redshift/mock"
	"github.com/HomesNZ/buyer-demand/internal/model"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"sort"
	"strings"
	"time"
)

var _ = Describe("BuyerDemand", func() {
	var mockCtrl *gomock.Controller
	var s Service
	var esClient *mockES.MockClient
	var rdClient *mockRD.MockClient

	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 1, 2, 3, 4, now.Location())
	lastMonth := time.Date(now.Year(), now.Month()-1, 1, 1, 2, 3, 4, now.Location())
	previousYear := time.Date(1990, now.Month(), 1, 1, 2, 3, 4, now.Location())
	mapItems := model.MapItemESs{
		// current listing: 4 Bedrooms, 3 Bathrooms, Halswell, RR1234, 850000
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(600000),
			ListingId:           null.StringFrom("listing-test-1"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		// current listing: 4 Bedrooms, 3 Bathrooms, Halswell, RR4321, 900000
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(900000),
			ListingId:           null.StringFrom("listing-test-2"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		// current listing: null Bedrooms, 3 Bathrooms, Halswell, RR1234, 900000
		model.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(650000),
			ListingId:           null.StringFrom("listing-test-3"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		// current listing: 4 Bedrooms, null Bathrooms, Halswell, RR1234, 900000
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.Int{},
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(850000),
			ListingId:           null.StringFrom("listing-test-4"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		// current listing: 4 Bedrooms, 3 Bathrooms, null Suburb, RR1234, 900000
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.String{},
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(725000),
			ListingId:           null.StringFrom("listing-test-5"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.TimeFrom(now.AddDate(-1, 0, -3)),
		},
		// current listing: 4 Bedrooms, 3 Bathrooms, Halswell, null CategoryCode, 900000
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.String{},
			Price:               null.FloatFrom(750000),
			ListingId:           null.StringFrom("listing-test-6"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.TimeFrom(now.AddDate(0, -2, -1)),
		},
		// current listing: 4 Bedrooms, 3 Bathrooms, Halswell, RP1234, null Price
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RP"),
			Price:               null.Float{},
			ListingId:           null.StringFrom("listing-test-7"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.TimeFrom(now.AddDate(0, 0, -10)),
		},
		// current listing: 4 Bedrooms, 3 Bathrooms, Halswell, RP1234, 900000
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RP"),
			Price:               null.FloatFrom(800000),
			ListingId:           null.StringFrom("listing-test-8"),
			LatestListingDate:   null.Time{},
			LatestSoldDate:      null.Time{},
		},
		// recent unlisted: null Bedrooms, 3 Bathrooms, Halswell, RR1234, 900000
		model.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(800000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(thisMonth),
			LatestSoldDate:      null.TimeFrom(now),
		},
		// recent unlisted: null Bedrooms, 3 Bathrooms, Halswell, RR1234, 900000
		model.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(830000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(lastMonth),
			LatestSoldDate:      null.TimeFrom(now),
		},
		// recent unlisted: null Bedrooms, 3 Bathrooms, Halswell, RR1234, 900000
		model.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(760000),
			ListingId:           null.String{},
			PropertyState:       null.Int{},
			LatestListingDate:   null.TimeFrom(previousYear),
			LatestSoldDate:      null.TimeFrom(now),
		},
		// recent unlisted: 4 Bedrooms, 3 Bathrooms, Halswell, RR1234, 900000
		model.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Suburb:              null.StringFrom("Halswell"),
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(730000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(lastMonth),
			LatestSoldDate:      null.TimeFrom(now),
		},
	}

	bds := model.BuyerDemands{
		model.BuyerDemand{
			NumBedrooms:            null.Int{},
			NumBathrooms:           null.IntFrom(3),
			Suburb:                 null.StringFrom("Halswell"),
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.IntFrom(45),
			MedianSalePrice:        650000,
			NumOfForSaleProperties: 1,
		},
		model.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.Int{},
			Suburb:                 null.StringFrom("Halswell"),
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        850000,
			NumOfForSaleProperties: 1,
		},
		model.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			Suburb:                 null.String{},
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        725000,
			NumOfForSaleProperties: 1,
		},
		model.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			Suburb:                 null.StringFrom("Halswell"),
			PropertyType:           null.String{},
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        750000,
			NumOfForSaleProperties: 1,
		},
		model.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			Suburb:                 null.StringFrom("Halswell"),
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.IntFrom(60),
			MedianSalePrice:        750000,
			NumOfForSaleProperties: 2,
		},
		model.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			Suburb:                 null.StringFrom("Halswell"),
			PropertyType:           null.StringFrom("RP"),
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        400000,
			NumOfForSaleProperties: 2,
		},
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		esClient = mockES.NewMockClient(mockCtrl)
		rdClient = mockRD.NewMockClient(mockCtrl)
		s = &service{
			esClient:       esClient,
			redshiftClient: rdClient,
			logger:         logrus.WithField("testing", true),
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("DailyBuyerDemandTableRefresh", func() {
		ctx := context.Background()

		It("query all from es error", func() {
			esClient.EXPECT().QueryAll(ctx).Return(nil, errors.New("test error"))

			err := s.DailyBuyerDemandTableRefresh(ctx)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("QueryAll"))
		})

		It("refresh table error", func() {
			esClient.EXPECT().QueryAll(ctx).Return(mapItems, nil)
			rdClient.EXPECT().DailyBuyerDemandTableRefresh(ctx, gomock.Any()).Return(errors.New("test error"))

			err := s.DailyBuyerDemandTableRefresh(ctx)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("DailyBuyerDemandTableRefresh"))
		})

		It("correct", func() {
			esClient.EXPECT().QueryAll(ctx).Return(mapItems, nil)
			rdClient.EXPECT().DailyBuyerDemandTableRefresh(ctx, gomock.Any()).Return(nil)

			err := s.DailyBuyerDemandTableRefresh(ctx)

			Expect(err).NotTo(HaveOccurred())
		})

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
