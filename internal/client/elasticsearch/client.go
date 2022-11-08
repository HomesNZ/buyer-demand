package elasticsearch

import (
	"context"
	"encoding/json"
	entity "github.com/HomesNZ/buyer-demand/internal/entity"
	"github.com/pkg/errors"

	"github.com/HomesNZ/elastic"
	"github.com/sirupsen/logrus"
)

var (
	// AliasName is the name of the ES alias which points to the index this schema resides in
	AliasName = "map_items"
)

type Client interface {
	BySuburbID(ctx context.Context, suburbID int) (entity.MapItemESs, error)
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

	// TODO Size From
	search := es.conn.Search().Index(AliasName).Type("map_item").Query(query).Size(10000)

	searchResult, err := search.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Do")
	}

	return parseToMapItem(searchResult.Hits.Hits)
}

func parseToMapItem(hits []*elastic.SearchHit) (entity.MapItemESs, error) {
	results := entity.MapItemESs{}
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
