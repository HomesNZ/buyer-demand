package service

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/HomesNZ/buyer-demand/internal/util"
	"github.com/pkg/errors"
	"sync"
	"time"
)

const suburbChunkSize = 100

func (s service) DailyBuyerDemandTableRefresh(ctx context.Context) error {
	// Query all suburb ids from db
	suburbIDs, err := s.repos.Address().AllSuburbIDs(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("AllSuburbIDs, %v", err))
		return errors.Wrap(err, "AllSuburbIDs")
	}

	// Query all properties/listings by suburb id from ES
	dataStartDate := util.ToUtcDate(time.Now().AddDate(0, 0, -360)).Format("2006-01-02T15:04:05Z")
	buyerDemands := make(chan entity.BuyerDemands, suburbChunkSize)
	needToDeleteTodayData := true
	wg := sync.WaitGroup{}
	for i, suburbID := range suburbIDs {
		if suburbID == 0 {
			continue
		}
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			// Calculate buyer demand, aggregated by num_bedrooms, num_bathrooms and property_sub_category
			bds, err := s.calculateBuyerDemands(ctx, id, dataStartDate)
			if err == nil && bds != nil && len(bds) > 0 {
				buyerDemands <- bds
			}
		}(suburbID)

		if (i+1)%suburbChunkSize == 0 {
			wg.Wait()
			s.logger.Info(fmt.Sprintf("Populate Batch #%d", (i+1)/suburbChunkSize))
			// Populate to DB
			for len(buyerDemands) > 0 {
				err = s.repos.BuyerDemand().Populate(ctx, <-buyerDemands, needToDeleteTodayData)
				if err != nil {
					s.logger.Error(fmt.Sprintf("BuyerDemand().Populate, %v", err))
					return errors.Wrap(err, "BuyerDemand().Populate")
				}
			}
			needToDeleteTodayData = false
		}
	}
	wg.Wait()
	close(buyerDemands)

	// Populate to DB
	s.logger.Info(fmt.Sprintf("Populate Last Batch"))
	for len(buyerDemands) > 0 {
		err = s.repos.BuyerDemand().Populate(ctx, <-buyerDemands, needToDeleteTodayData)
		if err != nil {
			s.logger.Error(fmt.Sprintf("BuyerDemand().Populate, %v", err))
			return errors.Wrap(err, "BuyerDemand().Populate")
		}
	}

	return nil
}

func (s service) calculateBuyerDemands(ctx context.Context, suburbID int, dataStartDate string) (entity.BuyerDemands, error) {
	mapItems, err := s.esClient.BySuburbID(ctx, suburbID, dataStartDate)
	if err != nil {
		s.logger.Error(fmt.Sprintf("esClient.BySuburbID (%d), %v", suburbID, err))
		return nil, errors.Wrap(err, fmt.Sprintf("esClient.BySuburbID %d", suburbID))
	}

	return mapItems.GenerateBuyerDemands(), nil
}
