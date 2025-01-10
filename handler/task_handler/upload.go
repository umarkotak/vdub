package task_handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

func UploadToYoutube(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	taskName := utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name"))

	taskDir := fmt.Sprintf("%s/%s", config.Get().BaseDir, taskName)
	state, err := service.GetState(ctx, taskDir, model.TaskState{})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	if state.Status != model.STATE_DUBBED_VIDEO_GENERATED {
		err = fmt.Errorf("need to finish dubbing for upload")
		utils.RenderError(w, r, 422, err)
		return
	}

	params := struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description"`
	}{}

	err = utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	if params.Description == "" {
		params.Description = "this is an auto dubbed video"
	}

	data, err := service.Upload(ctx, service.GenericUploadParams{
		VideoPath:   fmt.Sprintf("%s/dubbed_video.mp4", taskDir),
		Title:       params.Title,
		Description: params.Description,
	})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	utils.Render(w, r, 200, data, nil)
}
