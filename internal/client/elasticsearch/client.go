package elasticsearch

import (
	"context"
	"encoding/json"
	entity "github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/pkg/errors"

	"github.com/HomesNZ/elastic"
	"github.com/sirupsen/logrus"
)

const (
	// AliasName is the name of the ES alias which points to the index this schema resides in
	AliasName         = "map_items"
	bySuburbBatchSize = 1000
)

type Client interface {
	BySuburbID(ctx context.Context, suburbID int) (entity.MapItemESs, error)
	ByPropertyID(ctx context.Context, propertyID string) (*entity.MapItemES, error)
}

type client struct {
	log  *logrus.Entry
	conn *elastic.Client
}

func (es *client) BySuburbID(ctx context.Context, suburbID int) (entity.MapItemESs, error) {
	query := elastic.NewBoolQuery().Must(
		elastic.NewNestedQuery(
			"nested_address",
			elastic.NewBoolQuery().Must(
				elastic.NewMatchQuery("nested_address.suburb_id", suburbID))))

	return es.doSearch(ctx, query, 0, entity.MapItemESs{})
}

func (es *client) doSearch(ctx context.Context, query elastic.Query, from int, results entity.MapItemESs) (entity.MapItemESs, error) {
	search := es.conn.Search().Index(AliasName).Type("map_item").Query(query).Size(bySuburbBatchSize).From(from)
	searchResult, err := search.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Do")
	}

	results, err = parseToMapItem(results, searchResult.Hits.Hits)
	if err != nil {
		return nil, errors.Wrap(err, "parseToMapItem")
	}

	if searchResult.TotalHits() > int64(from+len(searchResult.Hits.Hits)) {
		return es.doSearch(ctx, query, from+len(searchResult.Hits.Hits), results)
	}

	return results, nil
}

func parseToMapItem(results entity.MapItemESs, hits []*elastic.SearchHit) (entity.MapItemESs, error) {
	for _, v := range hits {
		r := entity.MapItemES{}
		err := json.Unmarshal(*v.Source, &r)
		if err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}

		results = append(results, r)
	}
	return results, nil
}

func (es *client) ByPropertyID(ctx context.Context, propertyID string) (*entity.MapItemES, error) {
	query := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("property_id", propertyID))
	search := es.conn.Search().Index(AliasName).Type("map_item").Query(query).Size(1)
	searchResult, err := search.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Do")
	}

	if searchResult.TotalHits() < 1 {
		return nil, nil
	}

	results := entity.MapItemESs{}
	results, err = parseToMapItem(results, searchResult.Hits.Hits)
	if err != nil {
		return nil, errors.Wrap(err, "parseToMapItem")
	}

	return &results[0], nil
}
