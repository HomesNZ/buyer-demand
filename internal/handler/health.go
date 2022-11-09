package handler

import (
	"github.com/HomesNZ/buyer-demand/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Health(logger *logrus.Entry, s service.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.Health(); err != nil {
			logger.WithError(err).Error()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
