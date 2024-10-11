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

func GetTaskStatus(w http.ResponseWriter, r *http.Request) {
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

	utils.Render(
		w, r, 200,
		map[string]any{
			"state":       state,
			"state_human": state.GetTaskStateData(handlerState.RunningTask[taskName]),
		},
		nil,
	)
}

func GetTaskLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	taskName := utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name"))

	taskDir := fmt.Sprintf("%s/%s", config.Get().BaseDir, taskName)

	taskLog, err := utils.QuickGetLog(taskDir)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	utils.Render(w, r, 200, map[string]any{
		"logs": taskLog,
	}, nil,
	)
}

func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := struct {
		TaskName string `json:"-"`
		Status   string `json:"status"`
	}{
		TaskName: utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name")),
	}

	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	taskDir := utils.GenTaskDir(params.TaskName)

	state, err := service.GetState(ctx, taskDir, model.TaskState{})
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	err = service.SaveStateStatus(ctx, taskDir, &state, params.Status)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	utils.Render(
		w, r, 200,
		map[string]any{
			"state":       state,
			"state_human": state.GetTaskStateData(handlerState.RunningTask[params.TaskName]),
		},
		nil,
	)
}
