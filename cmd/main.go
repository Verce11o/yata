package main

import (
	"github.com/Verce11o/yata/internal/app"
	"github.com/Verce11o/yata/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	// TODO fix config loading after connect

	app.Run(cfg)
}
