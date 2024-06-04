package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

func GetTaskList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	myProjectPrefix := fmt.Sprintf("task-%s-", commonCtx.DirectUsername)

	files, err := os.ReadDir(config.Get().BaseDir)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	taskList := []model.TaskData{}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), myProjectPrefix) {
			taskName := file.Name()
			taskDir := fmt.Sprintf("%s/%s", config.Get().BaseDir, taskName)

			state, err := service.GetState(ctx, taskDir)
			stateDetail := state.GetTaskStateData(handlerState.RunningTask[taskName])
			if err != nil {
				logrus.WithContext(r.Context()).Error(err)
				utils.RenderError(w, r, 422, err)
				return
			}

			taskList = append(taskList, model.TaskData{
				Name:               strings.TrimPrefix(taskName, myProjectPrefix),
				Status:             stateDetail.Status,
				CurrentStatusHuman: stateDetail.CurrentStatusHuman,
				IsRunning:          stateDetail.IsRunning,
				ProgressSummary:    stateDetail.ProgressSummary,
			})
		}
	}

	utils.Render(w, r, 200, taskList, nil)
}
