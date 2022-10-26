package elasticsearch

import (
	"strings"

	"github.com/HomesNZ/elastic"
	"github.com/HomesNZ/go-common/elasticsearch"
	"github.com/HomesNZ/go-common/env"
	"github.com/sirupsen/logrus"
)

func New(log *logrus.Entry) (Client, error) {
	id := env.GetString("AWS_ACCESS_KEY_ID", "")
	secret := env.GetString("AWS_SECRET_ACCESS_KEY", "")
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(strings.Split(env.MustGetString("ELASTICSEARCH_URLS"), ";")...),
		elastic.SetHealthcheck(env.GetBool("ELASTICSEARCH_HEALTH_CHECK", true)),
		elastic.SetSniff(env.GetBool("ELASTICSEARCH_SNIFF", false)), // causes issues within AWS, so off by default
	}
	log.Info("Using AWS credentials for Elasticsearch")
	options = append(options, elasticsearch.AWSAccessKey(id, secret))
	esClient, err := elastic.NewClient(options...)
	if err != nil {
		log.Fatal(err)
	}
	return &client{log: log, conn: esClient}, nil
}
