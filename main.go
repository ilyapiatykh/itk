package main

import (
	"github.com/ilyapiatykh/itk/config"
	"github.com/ilyapiatykh/itk/internal/app"
)

func main() {
	cfg, err := config.NewCfg()
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
