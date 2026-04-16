package newsletter

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/dryrun"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	persistentdata "github.com/SeaweedbrainCY/jellyfin-newsletter/internal/persistentData"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/smtp"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/tmdb"
	"go.uber.org/zap"
)

type Workflow struct {
	JellyfinClient jellyfin.APIClient
	TMDBClient     tmdb.APIClient
}

// Run connects to Jellyfin to retrieve the latest items and send the newsletter to the configured recipients.
// cronjob is optional and should be nil if the workflow is not called by a scheduled job. It is mainly used for logging purposes
func (workflow Workflow) Run(app *app.ApplicationContext) {
	app.Logger.Info("Gathering new items and sending the newsletter ...")
	err := workflow.JellyfinClient.TestConnection(app)
	if err != nil {
		app.Logger.Fatal(
			"Jellyfin newsletter startup failed. An error occurred while connecting to Jellyfin.",
			zap.Error(err),
		)
	}

	recentlyAddedMovies := workflow.JellyfinClient.GetRecentlyAddedMovies(app)
	recentlyAddedSeries := workflow.JellyfinClient.GetNewlyAddedSeries(app)

	if len(*recentlyAddedMovies) == 0 && len(*recentlyAddedSeries) == 0 {
		app.Logger.Info("No new items detected. Email notification is skipped.")
		return
	}

	tmdb.EnrichMovieItemsList(recentlyAddedMovies, workflow.TMDBClient, app)
	tmdb.EnrichSeriesItemsList(recentlyAddedSeries, workflow.TMDBClient, app)

	moviesCount, episodesCount, err := workflow.JellyfinClient.LibraryAPI.GetItemsStats(app)
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
		err = smtp.SendEmailToAllRecipients(emailHTML, app)
		if err != nil {
			app.Logger.Fatal("Failed to send emails to recipients.", zap.Error(err))
		}
	}

	err = persistentdata.UpdateLastNewsletterDatetime(app.Clock.Now().UTC(), app)
	if err != nil {
		app.Logger.Warn(
			"An error occured while saving the last newsletter datetime. This could lead to future error or items sent again.",
			zap.Error(err),
		)
	}

	app.Logger.Info("Thanks for using Jellyfin-Newsletter !")
}
