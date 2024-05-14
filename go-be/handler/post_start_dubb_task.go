package handler

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

type (
	StartDubbTaskParams struct {
		TaskName   string `json:"task_name"`
		YoutubeUrl string `json:"youtube_url"`

		TaskDir      string
		RawVideoName string
		RawVideoPath string
	}
)

func (p *StartDubbTaskParams) Gen() {
	p.TaskDir = fmt.Sprintf("%s/%s", config.Get().BaseDir, p.TaskName)
	p.RawVideoName = "raw_video.mp4"
	p.RawVideoPath = fmt.Sprintf("%s/%s", p.TaskDir, p.RawVideoName)
}

func PostStartDubbTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := StartDubbTaskParams{}
	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}
	params.Gen()

	state, err := service.GetState(ctx, params.TaskDir)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	if state.Status == model.STATE_INITIALIZED {
		err = service.DownloadYoutubeVideo(ctx, params.YoutubeUrl, params.RawVideoPath)
		if err != nil {
			logrus.WithContext(r.Context()).Error(err)
			utils.RenderError(w, r, 422, err)
			return
		}

		err = service.SaveStateStatus(ctx, params.TaskDir, state, model.STATE_VIDEO_DOWNLOADED)
		if err != nil {
			logrus.WithContext(r.Context()).Error(err)
			utils.RenderError(w, r, 422, err)
			return
		}
	}
}
