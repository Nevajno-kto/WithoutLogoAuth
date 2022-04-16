package main

import (
	"github.com/nevajno-kto/without-logo-auth/config"
	"github.com/nevajno-kto/without-logo-auth/internal/app"
)

func main() {
	// Configuration
	cfg := config.GetConfig()
	// if err != nil {
	// 	log.Fatalf("Config error: %s", err)
	// }

	// Run
	app.Run(cfg)
}
