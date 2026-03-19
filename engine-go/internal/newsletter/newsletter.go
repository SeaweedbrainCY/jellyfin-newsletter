package newsletter

import (
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/dryrun"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	persistentdata "github.com/SeaweedbrainCY/jellyfin-newsletter/internal/persistentData"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/smtp"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/tmdb"
	"go.uber.org/zap"
)

// TriggerNewsletterWorkflow connects to Jellyfin to retrieve the latest items and send the newsletter to the configured recipients.
func TriggerNewsletterWorkflow(app *app.ApplicationContext) {
	app.Logger.Info("Gathering new items and sending the newsletter ...")
	jellyfinAPIClient := jellyfin.NewJellyfinAPIClient(app)
	err := jellyfinAPIClient.TestConnection(app)
	if err != nil {
		app.Logger.Fatal(
			"Jellyfin newsletter startup failed. An error occurred while connecting to Jellyfin.",
			zap.Error(err),
		)
	}

	recentlyAddedMovies := jellyfinAPIClient.GetRecentlyAddedMovies(app)
	recentlyAddedSeries := jellyfinAPIClient.GetNewlyAddedSeries(app)

	if len(*recentlyAddedMovies) == 0 && len(*recentlyAddedSeries) == 0 {
		app.Logger.Info("No new items detected. Email notification is skipped.")
		return
	}

	tmdbClient := tmdb.InitTMDBApiClient(app)
	tmdb.EnrichMovieItemsList(recentlyAddedMovies, tmdbClient, app)
	tmdb.EnrichSeriesItemsList(recentlyAddedSeries, tmdbClient, app)

	moviesCount, episodesCount, err := jellyfinAPIClient.LibraryAPI.GetItemsStats(app)
	if err != nil {
		app.Logger.Fatal("Failed to get Jellyfin items statistics.", zap.Error(err))
	}

	emailHTML, err := template.BuildNewMediaEmailHTML(
		recentlyAddedMovies,
		recentlyAddedSeries,
		moviesCount,
		episodesCount,
		app,
	)
	if err != nil {
		app.Logger.Fatal("Failed to build email HTML template.", zap.Error(err))
	}

	if app.Config.DryRun.Enabled {
		dryrun.SaveDryRunEmail(emailHTML, recentlyAddedMovies, recentlyAddedSeries, app)
		app.Logger.Info("Successfully generated the newsletter (dry run).")
	} else {
		smtp.SendEmailToAllRecipients(emailHTML, app)
	}

	err = persistentdata.UpdateLastNewsletterDatetime(time.Now(), app)
	if err != nil {
		app.Logger.Warn(
			"An error occured while saving the last newsletter datetime. This could lead to future error or items sent again.",
			zap.Error(err),
		)
	}

	app.Logger.Info("Thanks for using Jellyfin-Newsletter !")
}
