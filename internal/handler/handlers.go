package handler

import (
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/HomesNZ/go-common/env"
	"github.com/HomesNZ/go-common/logger"
	"github.com/HomesNZ/go-common/version"
	"github.com/HomesNZ/go-secret/auth"
	"github.com/HomesNZ/go-secret/auth/allow"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Register(log *logrus.Entry, r *mux.Router, a *auth.Auth, s service.Service) {
	cors := handlers.CORS(
		handlers.MaxAge(10000),
		handlers.AllowedMethods([]string{"GET"}),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedOrigins(strings.Split(env.GetString("CORS_ALLOWED_HOSTS", ""), ";")),
	)

	authorisation := &auth.Authorisation{
		Authenticator: a,
		Rules: auth.Rules{
			"buyer.demand.stats": allow.Any{
				allow.Role(auth.RoleUser),
			},
		},
	}

	stdChain := alice.New(
		logger.Middleware(log),
		cors,
		Gzip,
	)

	r.Handle("/version", stdChain.Then(handlers.MethodHandler{"GET": version.Handler}))
	r.Handle("/health", stdChain.Then(handlers.MethodHandler{"GET": Health(log.WithField("handler", "Health"), s)}))

	buyerDemandLatestStats := authorisation.MiddlewareAllow(
		"buyer.demand.stats",
		BuyerDemandLatestStats(log.WithField("handler", "BuyerDemandLatestStats"), s),
	)

	r.Handle("/stats/latest/{property_id}", stdChain.Then(handlers.MethodHandler{"GET": buyerDemandLatestStats}))

	http.Handle("/", r)
}
