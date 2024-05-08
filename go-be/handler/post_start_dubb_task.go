package handler

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/utils"
)

type (
	StartDubbTaskParams struct {
		TaskName   string `json:"task_name"`
		YoutubeUrl string `json:"youtube_url"`

		TaskDir string
	}
)

func (p *StartDubbTaskParams) Gen() {
	p.TaskDir = fmt.Sprintf("%s/%s", config.Get().BaseDir, p.TaskName)
}

func PostStartDubbTask(w http.ResponseWriter, r *http.Request) {
	params := StartDubbTaskParams{}
	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}
	params.Gen()

}
