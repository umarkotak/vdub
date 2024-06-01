package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

func GetTranscript(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	taskName := utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name"))

	transcriptInfo, err := service.GetTranscript(ctx, taskName, chi.URLParam(r, "transcript_type"))
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	utils.Render(
		w, r, 200,
		transcriptInfo,
		nil,
	)
}
