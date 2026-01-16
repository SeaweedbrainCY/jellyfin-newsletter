package main

import (
	"flag"
	"fmt"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/ctx"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/logger"
	"go.uber.org/zap"
)

var version = "dev" // Will be set during build time

func main() {
	var configPath = flag.String("config", "./config/config.yml", "path to config file")
	flag.Parse()

	config, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	logger, err := logger.LoadLogger(config)
	if err != nil {
		panic("an error occured while loading logger : " + err.Error())
	}

	context := ctx.InitContext(config, logger)

	context.Logger.Info("Starting Jellyfin Newsletter ...", zap.String("version", version))
	context.Logger.Info("Configuration loaded successfully")
}
