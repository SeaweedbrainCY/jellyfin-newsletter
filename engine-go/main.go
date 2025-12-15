package main

import (
	"flag"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"go.uber.org/zap"
)

var configPath = flag.String("config", "./config/config.yml", "path to config file")
var version = "dev" // Will be set during build time

func main() {
	flag.Parse()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Jellyfin Newsletter is starting ...", zap.String("version", version))
	logger.Info("Loading configuration", zap.String("configPath", *configPath))
	_, err := config.LoadConfiguration(*configPath, logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
		return
	}
	logger.Info("Configuration loaded successfully")

}
