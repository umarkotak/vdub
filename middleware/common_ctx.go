package middleware

import (
	"context"
	"net/http"

	"github.com/umarkotak/vdub-go/model"
)

func CommonContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		directUsername := "public"
		if r.Header.Get("x-direct-username") != "" {
			directUsername = r.Header.Get("x-direct-username")
		}

		commonCtx := model.CommonContext{
			DirectUsername: directUsername,
		}

		ctx := context.WithValue(r.Context(), "common_ctx", commonCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
