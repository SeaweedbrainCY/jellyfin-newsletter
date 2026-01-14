package main

import (
	"flag"
	"fmt"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/context"
	"go.uber.org/zap"
)

var version = "dev" // Will be set during build time

func main() {
	flag.Parse()
	var configPath = flag.String("config", "./config/config.yml", "path to config file")
	context, err := context.LoadContext(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	context.Logger.Info("Starting Jellyfin Newsletter ...", zap.String("version", version))
	context.Logger.Info("Configuration loaded successfully")
}
