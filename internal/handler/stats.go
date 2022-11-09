package handler

import (
	"github.com/HomesNZ/buyer-demand/internal/api"
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func BuyerDemandLatestStats(logger *logrus.Entry, s service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeSuburbAndCityRequest(r)
		if err != nil {
			EncodeErrorResponse(logger, w, err)
			return
		}

		res, err := s.BuyerDemandLatestStats(r.Context(), req)
		if err != nil {
			EncodeErrorResponse(logger, w, err)
			return
		}

		EncodeJSONResponse(logger, w, res)
	})
}

func decodeSuburbAndCityRequest(r *http.Request) (*api.BuyerDemandStatsRequest, error) {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	req := &api.BuyerDemandStatsRequest{}
	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		return nil, errors.Wrap(err, "schema.Decode")
	}
	return req, nil
}
