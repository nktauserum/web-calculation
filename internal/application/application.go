package application

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/veliashev/web-calculation/pkg/calculation"
)

type Application struct {
	Port int
}

type Request struct {
	Expression string `json:"expression"`
}

type Error struct {
	Error string `json:"error"`
}

type Response struct {
	Result float64 `json:"result"`
}

// Функция обрабатывает все ошибки, возвращая их в json-формате и с соответствующим кодом
func HandleError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	result := Error{Error: err.Error()}
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

func CalculationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Got a new request!")

	request := new(Request)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		HandleError(w, r, err, http.StatusUnprocessableEntity)
		return
	}

	response := Response{Result: result}
	bytes, err := json.Marshal(&response)
	if err != nil {
		HandleError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(bytes))

	log.Printf("Expression: %s", request.Expression)
	log.Printf("Result: %f", result)
}
func New(port int) *Application {
	return &Application{Port: port}
}

func (app *Application) Run() error {
	log.Println("Started!")
	http.HandleFunc("/api/v1/calculate", CalculationHandler)
	return http.ListenAndServe(":"+fmt.Sprint(app.Port), nil)
}
