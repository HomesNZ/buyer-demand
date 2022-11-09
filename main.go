package main

import (
	"context"
	es "github.com/HomesNZ/buyer-demand/internal/client/elasticsearch"
	"github.com/HomesNZ/buyer-demand/internal/handler"
	repo "github.com/HomesNZ/buyer-demand/internal/repository"
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/HomesNZ/go-common/bugsnag"
	"github.com/HomesNZ/go-common/env"
	"github.com/HomesNZ/go-common/logger"
	"github.com/HomesNZ/go-common/newrelic"
	"github.com/HomesNZ/go-common/version"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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
		log.Info("starting buyer demand cronjob")
		err := s.DailyBuyerDemandTableRefresh(ctx)
		if err != nil {
			bugsnag.Notify(err)
			log.WithError(err).Fatal()
		}
		log.Info("Buyer demand cronjob done")
		return
	}

	r := mux.NewRouter()
	handler.Register(log, r, s)
	addr := ":" + env.MustGetString("HTTP_PORT")
	log.Info("Listening on ", addr)
	logrus.Fatal(http.ListenAndServe(addr, r))
}
