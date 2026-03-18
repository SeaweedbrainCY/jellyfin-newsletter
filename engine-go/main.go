package main

import (
	"flag"
	"fmt"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/dryrun"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/i18n"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/logger"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/tmdb"
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

	localizer, err := i18n.NewLocalizer(config.EmailTemplate.Language)
	if err != nil {
		logger.Fatal("Failed to load Localizer", zap.Error(err))
	}

	app := app.InitApplicationContext(config, logger, localizer)

	err = template.CheckIfThemeIsAvailable(app)
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
		app.Logger.Fatal("Jellyfin newsletter startup failed. An error occurred while connecting to Jellyfin.", zap.Error(err))
	}

	recentlyAddedMovies := jellyfinAPIClient.GetRecentlyAddedMovies(app)
	recentlyAddedSeries := jellyfinAPIClient.GetNewlyAddedSeries(app)

	tmdbClient := tmdb.InitTMDBApiClient(app)
	tmdb.EnrichMovieItemsList(recentlyAddedMovies, tmdbClient, app)
	tmdb.EnrichSeriesItemsList(recentlyAddedSeries, tmdbClient, app)

	moviesCount, episodesCount, err := jellyfinAPIClient.LibraryAPI.GetItemsStats(app)
	if err != nil {
		app.Logger.Fatal("Failed to get Jellyfin items statistics.", zap.Error(err))
	}

	emailHTML, err := template.BuildNewMediaEmailHTML(recentlyAddedMovies, recentlyAddedSeries, moviesCount, episodesCount, app)
	if err != nil {
		app.Logger.Fatal("Failed to build email HTML template.", zap.Error(err))
	}

	if config.DryRun.Enabled {
		dryrun.SaveDryRunEmail(emailHTML, recentlyAddedMovies, recentlyAddedSeries, app)
		app.Logger.Info("Successfully generated the newsletter (dry run).")
	}

	app.Logger.Info("Thanks for using Jellyfin-Newsletter !")
	app.Logger.Info("Copyright (C) 2025 Nathan Stchepinsky (Seaweedbrain). Licensed under the AGPLv3.0")
}
