package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func BindJson(r *http.Request, dest any) error {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyByte, dest)
	if err != nil {
		return err
	}

	return nil
}
