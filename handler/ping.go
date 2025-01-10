package handler

import (
	"fmt"
	"net/http"

	"github.com/umarkotak/vdub-go/utils"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, 200, map[string]any{
		"ping": "pong",
	}, nil)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	utils.RenderError(w, r, 404, fmt.Errorf("route not found"))
}
