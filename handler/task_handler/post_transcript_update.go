package task_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/utils"
)

type (
	TranscriptUpdateParams struct {
		TaskName       string           `json:"task_name" validate:"required"` // must unique - it will determine the task folder
		TranscriptData []TranscriptData `json:"transcript_data"`
		// YoutubeUrl string `json:"youtube_url" validate:"required"` //
		// VoiceName      string `json:"voice_name" validate:"required"`  // eg: id-ID-ArdiNeural
		// VoiceRate      string `json:"voice_rate" validate:"required"`  // eg: [-/+]10%
		// VoicePitch     string `json:"voice_pitch" validate:"required"` // eg: [-/+]10Hz
		// ForceStartFrom string `json:"force_start_from"` // used to run from certain state
	}

	TranscriptData struct {
		StartAt string `json:"start_at"`
		EndAt   string `json:"end_at"`
		Value   string `json:"value"`
	}
)

func PostTranscriptUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := TranscriptUpdateParams{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
	}
	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}

	utils.Render(
		w, r, 200,
		nil,
		nil,
	)
}
