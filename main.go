package main

import (
	"log/slog"

	"github.com/ilyapiatykh/itk/config"
	"github.com/ilyapiatykh/itk/internal/app"
	"github.com/ilyapiatykh/itk/pkg/logging"
)

func main() {
	cfg, err := config.NewCfg()
	if err != nil {
		logging.Fatal("Failed to parse config", slog.Any("error", err))
	}

	app.Run(cfg)
}
