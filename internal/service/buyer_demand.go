package service

import (
	"context"
	"fmt"
	"github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const suburbChunkSize = 5

func (s service) DailyBuyerDemandTableRefresh(ctx context.Context) error {
	// Query all suburb ids from db
	suburbIDs, err := s.repos.Address().AllSuburbIDs(ctx)
	if err != nil {
		return errors.Wrap(err, "AllSuburbIDs")
	}

	// Query all properties/listings by suburb id from ES
	var buyerDemandsChan chan entity.BuyerDemands
	suburbIDChunks := chunkSlice(suburbIDs, suburbChunkSize)
	for _, chunk := range suburbIDChunks {
		g, c := errgroup.WithContext(ctx)
		for _, id := range chunk {
			id := id
			if id == 0 {
				continue
			}
			g.Go(func() error {
				// Calculate buyer demand, aggregate by num_bedrooms, num_bathrooms and property_sub_category
				bds, err := s.calculateBuyerDemands(c, id)
				if err != nil {
					return errors.Wrap(err, "calculateBuyerDemands")
				}
				if bds != nil {
					buyerDemandsChan <- bds
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return errors.Wrap(err, "g.Wait()")
		}
	}

	var buyerDemandsResult entity.BuyerDemands
	for bds := range buyerDemandsChan {
		buyerDemandsResult = append(buyerDemandsResult, bds...)
	}
	// Populate to DB
	err = s.repos.BuyerDemand().Populate(ctx, buyerDemandsResult)
	if err != nil {
		return errors.Wrap(err, "BuyerDemand().Populate")
	}

	s.logger.Info("DailyBuyerDemandTableRefresh is done")
	return nil
}

func chunkSlice(slice []int, chunkSize int) [][]int {
	var chunks [][]int
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

func (s service) calculateBuyerDemands(ctx context.Context, suburbID int) (entity.BuyerDemands, error) {
	mapItems, err := s.esClient.BySuburbID(ctx, suburbID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("esClient.BySuburbID %d", suburbID))
	}

	return mapItems.GenerateBuyerDemands(), nil
}
