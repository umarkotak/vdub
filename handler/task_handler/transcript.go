package task_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
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

func PatchTranscriptUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := model.TranscriptUpdateParams{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
	}
	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}

	// tmpTranscript, _ := json.MarshalIndent(params, " ", "  ")
	// logrus.Infof("TRANSCRIPT DATA: %+v\n", string(tmpTranscript))

	err = service.UpdateTranscript(ctx, params)
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

func PostTranscriptQuickShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := model.TranscriptUpdateParams{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
	}

	err := service.QuickShiftTranscriptV2(ctx, params)
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

func PostTranscriptDeleteByIdx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := model.TranscriptUpdatePosParams{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
		Idx:      utils.StringMustInt64(chi.URLParam(r, "idx")),
	}

	err := service.DeleteTranscriptByIdx(ctx, params)
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

func PostTranscriptAddNexyByIdx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := model.TranscriptUpdatePosParams{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
		Idx:      utils.StringMustInt64(chi.URLParam(r, "idx")),
	}

	err := service.AddTranscriptByIdx(ctx, params)
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

func PostGenPreviewVoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := model.TranscriptUpdatePosParams{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
		Idx:      utils.StringMustInt64(chi.URLParam(r, "idx")),
	}

	_, err := service.GenPreviewTranscriptVoice(ctx, params)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}

	utils.Render(w, r, 200, nil, nil)
}

func GetPreviewVoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := model.TranscriptUpdatePosParams{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
		Idx:      utils.StringMustInt64(chi.URLParam(r, "idx")),
	}

	previewPath, err := service.GetPreviewTranscriptVoice(ctx, params)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}

	http.ServeFile(w, r, previewPath)
}
