package jellyfin

import (
	"context"
	"strconv"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

type MovieItem struct {
	Id             string
	Name           string
	AdditionDate   time.Time
	TMDBId         int
	ProductionYear int
}

func (client *APIClient) GetRecentlyAddedMovies(app *app.ApplicationContext) *[]MovieItem {
	var movieItems = []MovieItem{}
	for _, folderName := range app.Config.Jellyfin.WatchedFilmFolders {
		if items, err := client.getRecentlyAddedMoviesByFolder(folderName, app); err == nil {
			movieItems = append(movieItems, items...)
		}
	}
	return &movieItems
}

func (client *APIClient) getRecentlyAddedMoviesByFolder(
	folderName string,
	app *app.ApplicationContext,
) ([]MovieItem, error) {
	minimumAdditionDate := time.Now().AddDate(0, 0, app.Config.Jellyfin.ObservedPeriodDays*-1)
	app.Logger.Debug(
		"Searching for recently added movies.",
		zap.String("FolderName", folderName),
		zap.String("StartAdditionDate", minimumAdditionDate.String()),
	)
	folderId, err := client.GetRootFolderIdByName(folderName, app)
	if err != nil {
		return nil, err
	}

	movies, getMoviesHTTPResponse, err := client.ItemsAPI.GetItems(context.Background()).
		Recursive(true).
		ParentId(folderId).
		LocationTypes([]api.LocationType{api.LOCATIONTYPE_FILE_SYSTEM}).
		IsMovie(true).
		Fields([]api.ItemFields{"DateCreated", "ProviderIds", "Id", "Name", "ProductionYear"}).
		Execute()

	if err != nil {
		logHttpResponseError(getMoviesHTTPResponse, err, app)
		return nil, err
	}

	var items = []MovieItem{}
	for _, movie := range movies.Items {
		name := "Unknown"
		TMDBId := 0
		productionYear := 0

		if movie.Name.IsSet() {
			name = *movie.Name.Get()
		}
		if movie.ProductionYear.IsSet() {
			productionYear = int(movie.GetProductionYear())
		}

		if !movie.DateCreated.IsSet() {
			app.Logger.Warn(
				"Movie ignored because it has no creation date.",
				zap.String("MovieID", *movie.Id),
				zap.String("MovieName", name),
			)
			continue
		}

		if value, ok := movie.ProviderIds["Tmdb"]; ok {
			if id, err := strconv.Atoi(value); err == nil {
				TMDBId = id
			}
		}

		if movie.DateCreated.Get().After(minimumAdditionDate) {
			items = append(items, MovieItem{
				Id:             *movie.Id,
				AdditionDate:   *movie.DateCreated.Get(),
				Name:           name,
				TMDBId:         TMDBId,
				ProductionYear: productionYear,
			})
		}
	}
	return items, nil
}
