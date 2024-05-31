package utils

import (
	"encoding/json"
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, statusCode int, data, err any) {
	if data == nil {
		data = map[string]any{}
	}

	if err == nil {
		err = map[string]any{}
	}

	res := map[string]any{
		"data":  data,
		"error": err,
	}
	b, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
}

func RenderError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	Render(w, r, statusCode, nil, map[string]any{
		"message": err.Error(),
	})
}

func SetCorsHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	w.Header().Add(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Animapu-User-Uid, X-Visitor-Id, X-From-Path",
	)
}
