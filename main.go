package main

import (
	"context"

	es "github.com/HomesNZ/buyer-demand/internal/client/elasticsearch"
	"github.com/HomesNZ/buyer-demand/internal/client/redshift"
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/HomesNZ/go-common/bugsnag"
	"github.com/HomesNZ/go-common/env"
	"github.com/HomesNZ/go-common/logger"
	"github.com/HomesNZ/go-common/newrelic"
	"github.com/HomesNZ/go-common/version"
	"github.com/HomesNZ/listing-service/service"
)

const ServiceName = "buyer-demand"

func main() {
	ctx := context.Background()
	env.InitEnv()
	// Initialize logger using go-common package + settings
	log := logger.Init(
		logger.Level(env.GetString("LOG_LEVEL", "info")),
	).WithField("service", ServiceName)
	version.Log()
	bugsnag.InitBugsnag()
	newrelic.InitNewRelic(ServiceName)

	elasticClient, err := es.New(log)
	if err != nil {
		bugsnag.Notify(err)
		log.WithError(err).Fatal()
	}

	redshiftClient, err := redshift.NewFromEnv(ctx)
	if err != nil {
		bugsnag.Notify(err)
		log.WithError(err).Fatal()
	}

	service, err := service.New(log, redshiftClient, elasticClient)
	if err != nil {
		bugsnag.Notify(err)
		log.WithError(err).Fatal()
	}

	if env.GetBool("BUYER_DEMAND_CRONJOB", false) {
		service.DailyBuyerDemandTableRefresh(ctx)
	}

}
