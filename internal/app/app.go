package app

import (
	"api-service/internal/config"
	"api-service/internal/http"
	taskservice "api-service/internal/service/task_service"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (*App) Run(cfg config.Config) {
	server := http.NewServer(cfg.ApiService.HTTP)
	taskService := taskservice.NewTaskService(cfg.Agent, cfg.DataProvider)
	taskHandler := taskservice.NewTaskHandler(taskService)
	taskHandler.Register(server.Mux)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		server.Run()
	}()
	defer func() {
		server.Stop()
	}()

	<-stop
}
