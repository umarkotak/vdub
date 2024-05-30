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
