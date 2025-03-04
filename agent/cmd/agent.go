package main

import (
	"log"

	"github.com/veliashev/web-calculation/agent/internal/controller"
)

func main() {
	agent := controller.NewAgent(8081)
	if err := agent.Run(); err != nil {
		log.Fatalf("error running agent: %s", err)
	}
}
