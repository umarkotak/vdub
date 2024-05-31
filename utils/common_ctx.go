package utils

import (
	"net/http"

	"github.com/umarkotak/vdub-go/model"
)

func GetCommonCtx(r *http.Request) model.CommonContext {
	v := r.Context().Value("common_ctx")

	if v == nil {
		return model.CommonContext{}
	}

	commonCtx, ok := v.(model.CommonContext)

	if !ok {
		return model.CommonContext{}
	}

	return commonCtx
}
