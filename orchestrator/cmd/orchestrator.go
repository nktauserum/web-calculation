package main

import (
	"log"

	"github.com/veliashev/web-calculation/orchestrator/internal/controller"
	"github.com/veliashev/web-calculation/shared/config"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error getting config: %s", err)
	}

	app := controller.New(config.Port)
	err = app.Run()
	if err != nil {
		log.Fatalf("Error starting application: %s", err)
	}
}
