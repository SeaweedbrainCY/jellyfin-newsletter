package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/i18n"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/logger"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"go.uber.org/zap"
)

//go:embed themes/*/*/*.html
var templateHTMLThemesFS embed.FS

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

	localizer, err := i18n.NewLocalizer(config.EmailTemplate.Language)
	if err != nil {
		logger.Fatal("Failed to load Localizer", zap.Error(err))
	}

	app := app.InitApplicationContext(config, logger, localizer)

	err = template.CheckIfThemeIsAvailable(templateHTMLThemesFS, app)
	if err != nil {
		app.Logger.Fatal(
			"Chosen theme doesn't exist or is not usable right now.",
			zap.String("Theme name", app.Config.EmailTemplate.Theme),
			zap.Error(err),
		)
	}

	app.Logger.Info("Starting Jellyfin Newsletter ...", zap.String("version", version))
	app.Logger.Info("Configuration loaded successfully")

	jellyfinAPIClient := jellyfin.NewJellyfinAPIClient(app)

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
