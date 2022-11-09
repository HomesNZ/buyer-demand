package handler

import (
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/HomesNZ/go-common/env"
	"github.com/HomesNZ/go-common/logger"
	"github.com/HomesNZ/go-common/newrelic"
	"github.com/HomesNZ/go-common/version"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Register(log *logrus.Entry, r *mux.Router, s service.Service) {
	cors := handlers.CORS(
		handlers.MaxAge(10000),
		handlers.AllowedMethods([]string{"GET", "PUT", "DELETE", "POST"}),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedOrigins(strings.Split(env.GetString("CORS_ALLOWED_HOSTS", ""), ";")),
	)

	stdChain := alice.New(
		newrelic.Middleware,
		logger.Middleware(log),
		cors,
		Gzip,
	)

	r.Handle("/version", stdChain.Then(handlers.MethodHandler{"GET": version.Handler}))
	r.Handle("/health", stdChain.Then(handlers.MethodHandler{"GET": Health(log.WithField("handler", "Health"), s)}))

	r.Handle("/stats/latest", stdChain.Then(handlers.MethodHandler{"GET": BuyerDemandLatestStats(log.WithField("handler", "BuyerDemandLatestStats"), s)}))

	http.Handle("/", r)
}
