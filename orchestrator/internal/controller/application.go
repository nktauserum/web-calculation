package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

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

	mux := mux.NewRouter()

	mux.HandleFunc("/api/v1/calculate", handler.CalculationHandler)
	mux.HandleFunc("/api/v1/expressions", handler.ExpressionsListHandler)
	mux.HandleFunc("/api/v1/expressions/{expressionID}", handler.ExpressionByIDHandler)
	mux.HandleFunc("/api/v1/tasks", handler.TaskListHandler)
	mux.HandleFunc("/internal/task", handler.GetAvailableTask)

	return http.ListenAndServe(":"+fmt.Sprint(app.Port), mux)
}
