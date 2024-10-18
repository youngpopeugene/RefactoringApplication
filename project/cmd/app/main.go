package main

import (
	"app/internal/app"
	"log"
)

func main() {
	err := app.Run()
	log.Fatal("App is stopped", "err", err)
}
