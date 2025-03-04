package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nktauserum/web-calculation/orchestrator/internal/service"
	tsk "github.com/nktauserum/web-calculation/orchestrator/pkg/task"
	"github.com/nktauserum/web-calculation/shared"
)

func GetAvailableTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queue := service.GetQueue()

	if r.Method == "POST" {
		var result shared.TaskResult
		body, err := io.ReadAll(r.Body)
		if err != nil {
			HandleError(w, r, err, http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &result)
		if err != nil {
			HandleError(w, r, err, http.StatusInternalServerError)
			return
		}

		queue.Done(result.ID, result.Result)
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Задача %d успешно выполнена!\n", result.ID)
		return
	}

	tasks := queue.GetTasks()
	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var finalTask *shared.Task
	var found bool

	for _, task := range tasks {
		if !task.Status && tsk.Complete(task) {
			finalTask = queue.FindTask(task.ID)
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(finalTask)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func ExpressionsListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queue := service.GetQueue()

	expressions := queue.GetExpressions()
	if len(expressions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(expressions)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func TaskListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queue := service.GetQueue()

	tasks := queue.GetTasks()
	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
