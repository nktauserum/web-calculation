package main

import (
	"log"

	"github.com/nktauserum/web-calculation/orchestrator/internal/controller"
)

func main() {
	app := controller.New()
	err := app.Run()
	if err != nil {
		log.Fatalf("Error starting application: %s", err)
	}
}
