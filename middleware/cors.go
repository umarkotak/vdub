package middleware

import (
	"net/http"

	"github.com/umarkotak/vdub-go/utils"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SetCorsHeaders(w)

		if r.Method == "OPTIONS" {
			utils.Render(w, r, 200, nil, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
