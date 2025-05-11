package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/nktauserum/web-calculation/orchestrator/internal/controller/middleware"
	"github.com/nktauserum/web-calculation/orchestrator/pkg/service"
	"github.com/nktauserum/web-calculation/shared"
)

func ExpressionsListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queue := service.GetQueue()
	userID := r.Context().Value(middleware.UserID).(int64)

	expressions := queue.GetExpressions()
	if len(expressions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var result []shared.Expression
	for _, expression := range expressions {
		if expression.UserID == userID {
			result = append(result, expression)
		}
	}
	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(result)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func ExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queue := service.GetQueue()

	vars := mux.Vars(r)
	expressionID, err := strconv.ParseInt(vars["expressionID"], 10, 64)
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	expressions := queue.GetExpressions()
	if len(expressions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID := r.Context().Value(middleware.UserID).(int64)
	currentTask := expressions[expressionID]

	if currentTask.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := json.Marshal(&currentTask)
	if err != nil {
		HandleError(w, r, err, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(resp)
}
