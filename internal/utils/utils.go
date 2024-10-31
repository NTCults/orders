package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"requestURI": r.RequestURI,
		}).Info()
		next.ServeHTTP(w, r)
	})
}
