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
	//var mockCtrl *gomock.Controller
	//var s Service
	//var esClient *mockES.MockClient
	//var repos *mockRepo.MockRepositories
	//var addressRepo *mockAddress.MockRepo
	//var buyerDemandRepo *mockBuyerDemand.MockRepo

	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 1, 2, 3, 4, now.Location())
	lastMonth := time.Date(now.Year(), now.Month()-1, 1, 1, 2, 3, 4, now.Location())
	previousYear := time.Date(1990, now.Month(), 1, 1, 2, 3, 4, now.Location())

	address1020 := entity.Address{
		SuburbID: null.IntFrom(1020),
	}
	addressNull := entity.Address{
		SuburbID: null.Int{},
	}
	mapItems := entity.MapItemESs{
		// current listing: 4 Bedrooms, 3 Bathrooms, 1020, RR1234, 850000
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
		// current listing: 4 Bedrooms, 3 Bathrooms, 1020, RR4321, 900000
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(900000),
			ListingId:           null.StringFrom("listing-test-2"),
			LatestListingDate:   null.TimeFrom(now),
			LatestSoldDate:      null.Time{},
		},
		// current listing: null Bedrooms, 3 Bathrooms, 1020, RR1234, 900000
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
		// current listing: 4 Bedrooms, null Bathrooms, 1020, RR1234, 900000
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
		// current listing: 4 Bedrooms, 3 Bathrooms, null SuburbID, RR1234, 900000
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
		// current listing: 4 Bedrooms, 3 Bathrooms, 1020, null CategoryCode, 900000
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
		// current listing: 4 Bedrooms, 3 Bathrooms, 1020, RP1234, null Price
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
		// current listing: 4 Bedrooms, 3 Bathrooms, 1020, RP1234, 900000
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
		// recent unlisted: null Bedrooms, 3 Bathrooms, 1020, RR1234, 900000
		entity.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(800000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(thisMonth),
			LatestSoldDate:      null.TimeFrom(now),
		},
		// recent unlisted: null Bedrooms, 3 Bathrooms, 1020, RR1234, 900000
		entity.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(830000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(lastMonth),
			LatestSoldDate:      null.TimeFrom(now),
		},
		// recent unlisted: null Bedrooms, 3 Bathrooms, 1020, RR1234, 900000
		entity.MapItemES{
			NumBedrooms:         null.Int{},
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(760000),
			ListingId:           null.String{},
			PropertyState:       null.Int{},
			LatestListingDate:   null.TimeFrom(previousYear),
			LatestSoldDate:      null.TimeFrom(now),
		},
		// recent unlisted: 4 Bedrooms, 3 Bathrooms, 1020, RR1234, 900000
		entity.MapItemES{
			NumBedrooms:         null.IntFrom(4),
			NumBathrooms:        null.IntFrom(3),
			Address:             address1020,
			PropertySubCategory: null.StringFrom("RR"),
			Price:               null.FloatFrom(730000),
			ListingId:           null.String{},
			PropertyState:       null.IntFrom(2),
			LatestListingDate:   null.TimeFrom(lastMonth),
			LatestSoldDate:      null.TimeFrom(now),
		},
	}

	bds := entity.BuyerDemands{
		entity.BuyerDemand{
			NumBedrooms:            null.Int{},
			NumBathrooms:           null.IntFrom(3),
			SuburbID:               null.IntFrom(1020),
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.IntFrom(21),
			MedianSalePrice:        650000,
			NumOfForSaleProperties: 1,
		},
		entity.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.Int{},
			SuburbID:               null.IntFrom(1020),
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        850000,
			NumOfForSaleProperties: 1,
		},
		entity.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			SuburbID:               null.Int{},
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        725000,
			NumOfForSaleProperties: 1,
		},
		entity.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			SuburbID:               null.IntFrom(1020),
			PropertyType:           null.String{},
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        750000,
			NumOfForSaleProperties: 1,
		},
		entity.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			SuburbID:               null.IntFrom(1020),
			PropertyType:           null.StringFrom("RR"),
			MedianDaysToSell:       null.IntFrom(37),
			MedianSalePrice:        750000,
			NumOfForSaleProperties: 2,
		},
		entity.BuyerDemand{
			NumBedrooms:            null.IntFrom(4),
			NumBathrooms:           null.IntFrom(3),
			SuburbID:               null.IntFrom(1020),
			PropertyType:           null.StringFrom("RP"),
			MedianDaysToSell:       null.Int{},
			MedianSalePrice:        400000,
			NumOfForSaleProperties: 2,
		},
	}

	BeforeEach(func() {
		//mockCtrl = gomock.NewController(GinkgoT())
		//esClient = mockES.NewMockClient(mockCtrl)
		//repos = mockRepo.NewMockRepositories(mockCtrl)
		//addressRepo = mockAddress.NewMockRepo(mockCtrl)
		//buyerDemandRepo = mockBuyerDemand.NewMockRepo(mockCtrl)
		//s = &service{
		//	repos:    repos,
		//	esClient: esClient,
		//	logger:   logrus.WithField("testing", true),
		//}
	})

	AfterEach(func() {
		//mockCtrl.Finish()
	})

	Describe("DailyBuyerDemandTableRefresh", func() {
		//ctx := context.Background()

		//It("query suburbIDs from DB error", func() {
		//	repos.EXPECT().Address().Return(addressRepo)
		//	addressRepo.EXPECT().AllSuburbIDs(ctx).Return(nil, errors.New("test error"))
		//
		//	err := s.DailyBuyerDemandTableRefresh(ctx)
		//
		//	Expect(err).To(HaveOccurred())
		//	Expect(err.Error()).To(ContainSubstring("AllSuburbIDs"))
		//})
		//
		//It("query mapItems from ES error", func() {
		//	repos.EXPECT().Address().Return(addressRepo)
		//	addressRepo.EXPECT().AllSuburbIDs(ctx).Return([]int{1020}, nil)
		//
		//	esClient.EXPECT().BySuburbID(ctx, 1020).Return(nil, errors.New("test error"))
		//
		//	err := s.DailyBuyerDemandTableRefresh(ctx)
		//
		//	Expect(err).To(HaveOccurred())
		//	Expect(err.Error()).To(ContainSubstring("esClient.BySuburbID"))
		//})
		//
		//It("correct", func() {
		//	repos.EXPECT().Address().Return(addressRepo)
		//	addressRepo.EXPECT().AllSuburbIDs(ctx).Return([]int{1020}, nil)
		//
		//	esClient.EXPECT().BySuburbID(ctx, 1020).Return(mapItems, nil)
		//
		//	err := s.DailyBuyerDemandTableRefresh(ctx)
		//
		//	Expect(err).NotTo(HaveOccurred())
		//})

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
