package main

import (
	"log"

	"github.com/Taya-san/portofolio-backend/config"
	"github.com/Taya-san/portofolio-backend/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
