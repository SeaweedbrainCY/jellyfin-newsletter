package main

import (
	"flag"
	"fmt"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
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

	app := app.InitApplicationContext(config, logger)

	app.Logger.Info("Starting Jellyfin Newsletter ...", zap.String("version", version))
	app.Logger.Info("Configuration loaded successfully")

	jellyfinAPIClient := jellyfin.GetJellyfinAPIClient(app)

	err = jellyfinAPIClient.TestConnection(app)
	if err != nil {
		app.Logger.Debug("An error occurred while connecting to Jellyfin.", zap.Error(err))
		app.Logger.Fatal("Jellyfin newsletter startup failed. Switch to debug to display Go error message")
	}

	recentlyAddedMovies := jellyfinAPIClient.GetRecentlyAddedMovies(app)
	app.Logger.Info("movies", zap.Any("movies", recentlyAddedMovies))
	recentlyAddedSeries := jellyfinAPIClient.GetNewlyAddedSeries(app)
	app.Logger.Info("series", zap.Any("series", recentlyAddedSeries))
}
