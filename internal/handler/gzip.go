package handler

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// Gzip is middleware that gzips HTTP responses
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gziphandler.GzipHandler(next).ServeHTTP(w, r)
	})
}
