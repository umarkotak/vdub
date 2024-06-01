package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/utils"
)

func ServeVideo(w http.ResponseWriter, r *http.Request) {
	commonCtx := utils.GetCommonCtx(r)

	taskName := utils.GenTaskName(commonCtx.DirectUsername, chi.URLParam(r, "task_name"))

	videoName := "raw_video.mp4"
	if chi.URLParam(r, "video_type") == "translated" {
		videoName = "dubbed_video.mp4"
	}

	mediaFile := fmt.Sprintf("%s/%s/%s", config.Get().BaseDir, taskName, videoName)

	// w.Header().Set("Content-Type", "application/x-mpegURL")
	http.ServeFile(w, r, mediaFile)
}
