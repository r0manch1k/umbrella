package main

import (
	"fmt"

	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("config error: %w", err))
	}

	app.Run(cfg)
}
