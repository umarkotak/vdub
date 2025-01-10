package handler

import (
	"fmt"
	"net/http"

	"github.com/bregydoc/gtranslate"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/utils"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	translated, err := gtranslate.TranslateWithParams(
		"hello my name is umar",
		gtranslate.TranslationParams{
			From: "en", To: "id",
		},
	)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
	}

	utils.Render(w, r, 200, map[string]any{
		"ping":       "pong",
		"translated": translated,
	}, nil)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	utils.RenderError(w, r, 404, fmt.Errorf("route not found"))
}
