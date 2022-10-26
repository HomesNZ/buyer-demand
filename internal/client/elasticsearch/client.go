package elasticsearch

import (
	"context"

	"github.com/HomesNZ/elastic"
	"github.com/sirupsen/logrus"
)

var (
	// AliasName is the name of the ES alias which points to the index this schema resides in
	AliasName = "map_items"
)

type Client interface {
	QueryAllListings(ctx context.Context) (*elastic.SearchResult, error)
}

type client struct {
	log  *logrus.Entry
	conn *elastic.Client
}

func (es *client) QueryAllListings(ctx context.Context) (*elastic.SearchResult, error) {

	query := elastic.NewBoolQuery().Filter(
		elastic.NewExistsQuery("listing_id"),
	)
	search := es.conn.
		Search().
		Index(AliasName).
		Type("map_item").
		Query(query)

	return search.Do(ctx)
}
