package controller

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/nktauserum/web-calculation/orchestrator/internal/controller/handler"
	"github.com/nktauserum/web-calculation/orchestrator/internal/controller/middleware"
	"github.com/nktauserum/web-calculation/orchestrator/pkg/auth"
	"github.com/nktauserum/web-calculation/proto"
	"github.com/nktauserum/web-calculation/proto/pb"
)

type Orchestrator struct {
	Port        int
	DBPath      string
	JWTSecret   string
	TokenExpiry time.Duration
	grpc        *RPCServer
}

type RPCServer struct {
	server *proto.Server
	port   int
}

func NewRPCServer(port int) *RPCServer {
	return &RPCServer{
		port: port,
	}
}

func (s *RPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	s.server = &proto.Server{}
	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, s.server)

	log.Printf("Starting gRPC server on port %d", s.port)
	return grpcServer.Serve(lis)
}

func New() *Orchestrator {
	_ = godotenv.Load(".env")
	port := 8080 // Default port
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &Orchestrator{
		Port:        port,
		DBPath:      os.Getenv("DB_PATH"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		TokenExpiry: 24 * time.Hour, // Токен действителен 24 часа
		grpc:        NewRPCServer(5000),
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

	// запускаем gRPC сервер
	go func() {
		if err := app.grpc.Start(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	router := mux.NewRouter()

	// Публичные маршруты (без авторизации)
	router.HandleFunc("/api/v1/auth/register", handler.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", handler.LoginHandler).Methods("POST")

	// Защищенные маршруты (требуют авторизации)
	router.HandleFunc("/api/v1/calculate", authMiddleware.RequireAuth(handler.CalculationHandler))
	router.HandleFunc("/api/v1/expressions", authMiddleware.RequireAuth(handler.ExpressionsListHandler))
	router.HandleFunc("/api/v1/expressions/{expressionID}", authMiddleware.RequireAuth(handler.ExpressionByIDHandler))

	return http.ListenAndServe(":"+fmt.Sprint(app.Port), router)
}
