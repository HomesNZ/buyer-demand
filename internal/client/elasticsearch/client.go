package elasticsearch

import (
	"context"
	"encoding/json"
	mapItemES "github.com/HomesNZ/buyer-demand/internal/model"
	"github.com/pkg/errors"

	"github.com/HomesNZ/elastic"
	"github.com/sirupsen/logrus"
)

var (
	// AliasName is the name of the ES alias which points to the index this schema resides in
	AliasName = "map_items"
)

type Client interface {
	QueryAll(ctx context.Context) (mapItemES.MapItemESs, error)
}

type client struct {
	log  *logrus.Entry
	conn *elastic.Client
}

func (es *client) QueryAll(ctx context.Context) (mapItemES.MapItemESs, error) {
	query := elastic.NewBoolQuery().Filter(
		elastic.NewExistsQuery("property_id"),
	)
	search := es.conn.
		Search().
		Index(AliasName).
		Type("map_item").
		Query(query)

	searchResult, err := search.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Do")
	}

	return parseToMapItem(searchResult.Hits.Hits)
}

func parseToMapItem(hits []*elastic.SearchHit) (mapItemES.MapItemESs, error) {
	results := mapItemES.MapItemESs{}
	for _, v := range hits {
		r := mapItemES.MapItemES{}
		err := json.Unmarshal(*v.Source, &r)
		if err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}

		results = append(results, r)
	}
	return results, nil
}
