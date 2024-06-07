package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/datastore"
	"github.com/umarkotak/vdub-go/handler"
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
	)

	handler.Initialize()

	r.Get("/", handler.Ping)

	r.Post("/vdub/api/dubb/start", handler.PostStartDubbTask)
	r.Delete("/vdub/api/dubb/task/{task_name}", handler.DeleteTask)
	r.Get("/vdub/api/dubb/tasks", handler.GetTaskList)
	r.Get("/vdub/api/dubb/task/{task_name}/status", handler.GetTaskStatus)
	r.Get("/vdub/api/dubb/task/{task_name}/transcript/{transcript_type}", handler.GetTranscript)
	r.Post("/vdub/api/dubb/task/{task_name}/transcript/update", handler.PostTranscriptUpdate)

	r.Get("/vdub/api/dubb/task/{task_name}/video/{video_type}", handler.ServeVideo)
	r.Get("/vdub/api/dubb/task/{task_name}/video/snapshot", handler.ServeSnapshot)

	port := ":29000"
	logrus.Infof("Listening on port %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		logrus.Fatal(err)
	}
}
