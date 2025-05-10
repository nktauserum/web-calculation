package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nktauserum/web-calculation/orchestrator/pkg/service"
	"github.com/nktauserum/web-calculation/shared"
)

func CalculationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Got a new request!")

	query := new(shared.ExpressionRequest)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, query)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	queue := service.GetQueue()
	exprID, err := queue.ParseExpression(query.Expression)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	resp := shared.CalculateResponse{
		ID: exprID,
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
