package main

import (
	"log"

	"github.com/dubovayaRoshcha/posts-service/config"
	"github.com/dubovayaRoshcha/posts-service/internal/app"
)

func main() {
	err := config.Read()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(&config.Config)
}
