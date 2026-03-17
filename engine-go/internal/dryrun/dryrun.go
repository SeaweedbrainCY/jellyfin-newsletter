package dryrun

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
)

func SaveDryRunEmail(emailHTML string, newJellyfinMovies *[]jellyfin.MovieItem,
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem, app *app.ApplicationContext) {

		outputDirectory := "/app/config/previews/"
		if app.Config.DryRun.OutputDirectory == "" {
			outputDirectory = app.Config.DryRun.OutputDirectory
		}

		if app.Config.DryRun.OutputFilename := "newsletter_{date}.html"
}
