package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/clock"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/cron"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/i18n"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/logger"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/newsletter"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/tmdb"
	"github.com/go-co-op/gocron/v2"
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

	app := app.InitApplicationContext(config, logger, localizer, clock.RealClock{})

	err = template.CheckIfThemeIsAvailable(app)
	if err != nil {
		app.Logger.Fatal(
			"Chosen theme doesn't exist or is not usable right now.",
			zap.String("Theme name", app.Config.EmailTemplate.Theme),
			zap.Error(err),
		)
	}

	app.Logger.Info("Starting Jellyfin Newsletter ...", zap.String("version", version))
	app.Logger.Info("Copyright (C) 2025 Nathan Stchepinsky (Seaweedbrain). Licensed under the AGPLv3.0")
	app.Logger.Info("Configuration loaded successfully")

	newsletterWorkflow := newsletter.Workflow{
		JellyfinClient: jellyfin.NewJellyfinAPIClient(http.DefaultClient, app),
		TMDBClient:     tmdb.InitTMDBApiClient(http.DefaultClient, app),
	}

	if app.Config.Scheduler.Enabled {
		var scheduler gocron.Scheduler
		scheduler, err = cron.CreateNewsletterScheduler(newsletterWorkflow, app)
		if err != nil {
			app.Logger.Fatal("Error while creating the scheduler. Exiting now.", zap.Error(err))
		}
		scheduler.Start()
		// Block forever (or until signal)
		select {}
	}

	// One time trigger
	newsletterWorkflow.Run(app)

	app.Logger.Info("Jellyfin-Newsletter exiting gracefully.")
}
