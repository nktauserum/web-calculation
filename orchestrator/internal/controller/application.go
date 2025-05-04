package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/nktauserum/web-calculation/orchestrator/internal/controller/handler"
	"github.com/nktauserum/web-calculation/orchestrator/internal/controller/middleware"
	"github.com/nktauserum/web-calculation/orchestrator/pkg/auth"
)

type Orchestrator struct {
	Port        int
	DBPath      string
	JWTSecret   string
	TokenExpiry time.Duration
}

func New(port int, dbPath string, jwtSecret string) *Orchestrator {
	return &Orchestrator{
		Port:        port,
		DBPath:      dbPath,
		JWTSecret:   jwtSecret,
		TokenExpiry: 24 * time.Hour, // Токен действителен 24 часа
	}
}

func (app *Orchestrator) Run() error {
	log.Println("Orchestrator started!")

	userStorage, err := auth.NewUserStorage(app.DBPath)
	if err != nil {
		return fmt.Errorf("ошибка инициализации хранилища пользователей: %w", err)
	}
	defer userStorage.Close()

	authService := auth.NewAuthService(userStorage, app.JWTSecret, app.TokenExpiry)
	handler.SetAuthService(authService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	router := mux.NewRouter()

	// Публичные маршруты (без авторизации)
	router.HandleFunc("/api/v1/auth/register", handler.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", handler.LoginHandler).Methods("POST")

	// Защищенные маршруты (требуют авторизации)
	router.HandleFunc("/api/v1/calculate", authMiddleware.RequireAuth(handler.CalculationHandler))
	router.HandleFunc("/api/v1/expressions", authMiddleware.RequireAuth(handler.ExpressionsListHandler))
	router.HandleFunc("/api/v1/expressions/{expressionID}", authMiddleware.RequireAuth(handler.ExpressionByIDHandler))
	router.HandleFunc("/api/v1/tasks", authMiddleware.RequireAuth(handler.TaskListHandler))

	// Внутренние маршруты (не требуют авторизации)
	router.HandleFunc("/internal/task", handler.GetAvailableTask)

	return http.ListenAndServe(":"+fmt.Sprint(app.Port), router)
}
