package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/utils"
)

type (
	TaskData struct {
		Name string `json:"name"`
	}
)

func GetTaskList(w http.ResponseWriter, r *http.Request) {
	commonCtx := utils.GetCommonCtx(r)

	myProjectPrefix := fmt.Sprintf("task-%s", commonCtx.DirectUsername)

	files, err := os.ReadDir(config.Get().BaseDir)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	taskList := []TaskData{}

	for _, file := range files {
		if file.IsDir() {
			if strings.HasPrefix(file.Name(), myProjectPrefix) {
				taskList = append(taskList, TaskData{
					Name: file.Name(),
				})
			}
		}
	}

	utils.Render(
		w, r, 200,
		taskList,
		nil,
	)
}
