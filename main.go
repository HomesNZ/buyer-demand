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
	"github.com/HomesNZ/go-secret/auth"
	"github.com/HomesNZ/go-secret/auth/allow"
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

	a, err := auth.New(
		auth.Logger(log),
		auth.JWKS(env.MustGetString("AUTH0_JWKS_URL")),
		auth.ClientSecret(env.MustGetString("AUTH0_CLIENT_SECRET")),
		auth.APISecret(env.MustGetString("AUTH0_API_SECRET")),
		auth.ServiceKey(env.MustGetString("AUTHORISE_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}

	authorisation := &auth.Authorisation{
		Authenticator: a,
		Rules: auth.Rules{
			"buyer.demand.stats": allow.Always{},
		},
	}
	r := mux.NewRouter()
	handler.Register(log, r, authorisation, s)
	addr := ":" + env.MustGetString("HTTP_PORT")
	log.Info("Listening on ", addr)
	logrus.Fatal(http.ListenAndServe(addr, r))
}
