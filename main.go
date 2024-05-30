package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/datastore"
	"github.com/umarkotak/vdub-go/handler"
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
		chiMiddleware.RealIP,
		chiMiddleware.Recoverer,
	)

	handler.Initialize()

	r.Get("/", handler.Ping)

	r.Post("/vdub/api/dubb/start", handler.PostStartDubbTask)
	r.Get("/vdub/api/dubb/{task_name}/status", handler.GetTaskStatus)

	port := ":29000"
	logrus.Infof("Listening on port %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		logrus.Fatal(err)
	}
}
