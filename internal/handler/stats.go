package handler

import (
	"github.com/HomesNZ/buyer-demand/internal/api"
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/HomesNZ/buyer-demand/internal/util"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func BuyerDemandLatestStats(logger *logrus.Entry, s service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeBuyerDemandLatestStatsRequest(r)
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

func decodeBuyerDemandLatestStatsRequest(r *http.Request) (*api.BuyerDemandLatestStatsRequest, error) {
	vars := mux.Vars(r)
	propertyID := vars["property_id"]
	req := api.BuyerDemandLatestStatsRequest{
		PropertyID: propertyID,
	}

	decoder := schema.NewDecoder()
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		return nil, errors.Wrap(err, "decoder.Decode")
	}

	err = validation.ValidateStruct(&req,
		validation.Field(&req.PropertyID, validation.Required, is.UUID),
	)
	if err != nil {
		return nil, util.BadRequest(err.Error())
	}

	return &req, nil
}
