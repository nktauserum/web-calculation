package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Функция обрабатывает все ошибки, возвращая их в json-формате и с соответствующим кодом
func HandleError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	result := struct {
		Error string `json:"result"`
	}{Error: err.Error()}

	bytes, err := json.Marshal(&result)
	if err != nil {
		// Если произошли проблемы на этом моменте - не думая возвращаем 500
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unexpected error: %s", err)
		return
	}
	w.WriteHeader(statusCode)
	fmt.Fprint(w, string(bytes))
	log.Printf("Error: %s", result.Error)
	log.Printf("Status code: %d", statusCode)
}
