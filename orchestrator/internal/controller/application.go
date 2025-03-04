package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nktauserum/web-calculation/orchestrator/internal/controller/handler"
)

type Orchestrator struct {
	Port int
}

func New(port int) *Orchestrator {
	return &Orchestrator{Port: port}
}

func (app *Orchestrator) Run() error {
	log.Println("Orchestrator started!")

	http.HandleFunc("/api/v1/calculate", handler.CalculationHandler)
	http.HandleFunc("/api/v1/expressions", handler.ExpressionsListHandler)
	http.HandleFunc("/api/v1/tasks", handler.TaskListHandler)
	http.HandleFunc("/internal/task", handler.GetAvailableTask)

	return http.ListenAndServe(":"+fmt.Sprint(app.Port), nil)
}
