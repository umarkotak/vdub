package handler

import (
	"net/http"

	"github.com/umarkotak/vdub-go/utils"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, 200, map[string]any{
		"ping": "pong",
	}, nil)
}
