package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithContext(r.Context()).Infof(r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
