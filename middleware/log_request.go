package middleware

import (
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement logic

		next.ServeHTTP(w, r)
	})
}
