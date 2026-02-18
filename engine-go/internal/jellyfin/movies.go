package jellyfin

import (
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"go.uber.org/zap"
)

type MovieItem struct {
	ID             string
	Name           string
	AdditionDate   *time.Time
	TMDBId         string
	ProductionYear int32
}

// GetRecentlyAddedMovies aggregates recently added movies from all
// configured watched film folders. It queries each folder via
// `getRecentlyAddedMoviesByFolder` and returns a slice of
// `MovieItem` for movies added within the configured observed period.
func (client *APIClient) GetRecentlyAddedMovies(app *app.ApplicationContext) *[]MovieItem {
	var movieItems = []MovieItem{}
	for _, folderName := range app.Config.Jellyfin.WatchedFilmFolders {
		if items, err := client.getRecentlyAddedMoviesByFolder(folderName, app); err == nil {
			movieItems = append(movieItems, items...)
		}
	}
	return &movieItems
}

// getRecentlyAddedMoviesByFolder retrieves movies for a single
// Jellyfin folder, filters them by their addition/creation date
// against the configured observed period, and returns only those
// movies added after the computed cutoff. Movies without a
// creation date are logged and ignored.
func (client *APIClient) getRecentlyAddedMoviesByFolder(
	folderName string,
	app *app.ApplicationContext,
) ([]MovieItem, error) {
	minimumAdditionDate := time.Now().AddDate(0, 0, app.Config.Jellyfin.ObservedPeriodDays*-1-1)
	app.Logger.Debug(
		"Searching for recently added movies.",
		zap.String("FolderName", folderName),
		zap.String("StartAdditionDate", minimumAdditionDate.String()),
	)
	folderID, err := client.ItemsAPI.GetRootFolderIDByName(folderName, app)
	if err != nil {
		return nil, err
	}

	movies, err := client.ItemsAPI.GetMoviesItemsByFolderID(folderID, true, app)

	if err != nil {
		return nil, err
	}

	var items = []MovieItem{}
	for _, movie := range *movies {
		name := OrDefault(movie.Name, "")
		productionYear := OrDefault(movie.ProductionYear, 0)

		if !movie.DateCreated.IsSet() || movie.DateCreated.Get() == nil {
			app.Logger.Warn(
				"Movie ignored because it has no creation date.",
				zap.String("MovieID", *movie.Id),
				zap.String("MovieName", name),
			)
			continue
		}

		tmdbID := getTMDBIDIfExist(&movie)
		if movie.DateCreated.Get().After(minimumAdditionDate) {
			items = append(items, MovieItem{
				ID:             *movie.Id,
				AdditionDate:   movie.DateCreated.Get(),
				Name:           name,
				TMDBId:         tmdbID,
				ProductionYear: productionYear,
			})
		}
	}
	return items, nil
}
