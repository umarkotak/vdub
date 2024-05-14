package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

func GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskDir := fmt.Sprintf("%s/%s", config.Get().BaseDir, chi.URLParam(r, "task_name"))
	state, err := service.GetState(ctx, taskDir)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	utils.Render(w, r, 200, state, nil)
}
