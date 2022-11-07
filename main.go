package main

import (
	"context"
	es "github.com/HomesNZ/buyer-demand/internal/client/elasticsearch"
	repo "github.com/HomesNZ/buyer-demand/internal/repository"
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/HomesNZ/go-common/bugsnag"
	"github.com/HomesNZ/go-common/env"
	"github.com/HomesNZ/go-common/logger"
	"github.com/HomesNZ/go-common/newrelic"
	"github.com/HomesNZ/go-common/version"
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

	repos, err := repo.New(ctx)
	if err != nil {
		bugsnag.Notify(err)
		log.WithError(err).Fatal()
	}

	s, err := service.New(log, repos, elasticClient)
	if err != nil {
		bugsnag.Notify(err)
		log.WithError(err).Fatal()
	}

	if env.GetBool("BUYER_DEMAND_CRONJOB", false) {
		err := s.DailyBuyerDemandTableRefresh(ctx)
		if err != nil {
			bugsnag.Notify(err)
			log.WithError(err).Fatal()
		}
	}

}
