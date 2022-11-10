package service

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/pkg/errors"
	"sync"
)

const suburbChunkSize = 20

func (s service) DailyBuyerDemandTableRefresh(ctx context.Context) error {
	// Query all suburb ids from db
	suburbIDs, err := s.repos.Address().AllSuburbIDs(ctx)
	if err != nil {
		return errors.Wrap(err, "AllSuburbIDs")
	}

	// Query all properties/listings by suburb id from ES
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
			// Calculate buyer demand, aggregate by num_bedrooms, num_bathrooms and property_sub_category
			bds, err := s.calculateBuyerDemands(ctx, id)
			if err == nil && bds != nil {
				buyerDemands <- bds
			}
		}(suburbID)

		if (i+1)%suburbChunkSize == 0 {
			wg.Wait()

			// Populate to DB
			for len(buyerDemands) > 0 {
				err = s.repos.BuyerDemand().Populate(ctx, <-buyerDemands, needToDeleteTodayData)
				if err != nil {
					return errors.Wrap(err, "BuyerDemand().Populate")
				}
			}
			needToDeleteTodayData = false
		}
	}
	wg.Wait()
	close(buyerDemands)

	// Populate to DB
	for len(buyerDemands) > 0 {
		err = s.repos.BuyerDemand().Populate(ctx, <-buyerDemands, needToDeleteTodayData)
		if err != nil {
			return errors.Wrap(err, "BuyerDemand().Populate")
		}
	}

	s.logger.Info("DailyBuyerDemandTableRefresh is done")
	return nil
}

func (s service) calculateBuyerDemands(ctx context.Context, suburbID int) (entity.BuyerDemands, error) {
	mapItems, err := s.esClient.BySuburbID(ctx, suburbID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("esClient.BySuburbID %d", suburbID))
	}

	return mapItems.GenerateBuyerDemands(), nil
}
