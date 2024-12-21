package main

import (
	"log"

	"github.com/veliashev/rpn/internal/application"
)

func main() {
	app := application.New(8080)
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
