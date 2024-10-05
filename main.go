package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/datastore"
	"github.com/umarkotak/vdub-go/handler"
	"github.com/umarkotak/vdub-go/handler/task_handler"
	"github.com/umarkotak/vdub-go/middleware"
)

func initialize() {
	logrus.SetReportCaller(true)
	config.InitConfig()
	datastore.InitDataStore()
}

func main() {
	initialize()

	r := chi.NewRouter()

	r.Use(
		chiMiddleware.RequestID,
		chiMiddleware.Recoverer,
		middleware.CommonContext,
		middleware.Cors,
		middleware.LogRequest,
	)
	r.NotFound(handler.NotFound)

	task_handler.Initialize()

	r.Get("/", handler.Ping)

	r.Post("/vdub/api/dubb/start", task_handler.PostStartDubbTask)
	r.Delete("/vdub/api/dubb/task/{task_name}", task_handler.DeleteTask)
	r.Get("/vdub/api/dubb/tasks", task_handler.GetTaskList)
	r.Get("/vdub/api/dubb/task/{task_name}/status", task_handler.GetTaskStatus)
	r.Patch("/vdub/api/dubb/task/{task_name}/status", task_handler.UpdateTaskStatus)
	r.Patch("/vdub/api/dubb/task/{task_name}/transcript", task_handler.PatchTranscriptUpdate)
	r.Get("/vdub/api/dubb/task/{task_name}/transcript/{transcript_type}", task_handler.GetTranscript)

	r.Get("/vdub/api/dubb/task/{task_name}/video/{video_type}", task_handler.ServeVideo)
	r.Get("/vdub/api/dubb/task/{task_name}/video/snapshot", task_handler.ServeSnapshot)

	port := ":29000"
	logrus.Infof("Listening on port %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		logrus.Fatal(err)
	}
}
