package main

import (
	"log"

	"github.com/nktauserum/web-calculation/agent/internal/controller"
)

func main() {
	agent, err := controller.NewAgent("localhost:5000")
	if err != nil {
		log.Fatal(err)
	}

	if err := agent.Run(); err != nil {
		log.Fatalf("error running agent: %s", err)
	}
}
