package main

import (
	"log"

	"github.com/ttymayor/url-shortener/internal/app"
)

func main() {
	server := app.NewApp()
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
