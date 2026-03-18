package tmdb

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"go.uber.org/zap"
)

func enrichMovieWithDefaultInfos(jellyfinSeriesItem *jellyfin.MovieItem) {
	defaultDetails := getDefaultItemDetails()
	jellyfinSeriesItem.Overview = defaultDetails.Overview
	jellyfinSeriesItem.PosterURL = defaultDetails.PosterURL
}

func EnrichMovieItem(
	jellyfinMovieItem *jellyfin.MovieItem,
	tmdbAPIClient APIInterface,
	app *app.ApplicationContext,
) {
	if jellyfinMovieItem.TMDBId != "" {
		parsedHTTPResponse, err := tmdbAPIClient.GetMediaByID(jellyfinMovieItem.TMDBId, MediaTypeMovie)

		if err != nil {
			// Error is already logged by GetMediaByID
			enrichMovieWithDefaultInfos(jellyfinMovieItem)
			return
		}

		details := getItemDetailsFromHTTPResponse(parsedHTTPResponse)
		jellyfinMovieItem.Overview = details.Overview
		jellyfinMovieItem.PosterURL = details.PosterURL
		return
	}
	// No TMDB id, we perform a search by name and select the item with the highest popularity
	app.Logger.Debug(
		"Movie has no TMDB id. TMDB information will be retrieved by searching with Movie's name. If several media match, the choice will be based on the highest popularity.",
		zap.String("Series Name", jellyfinMovieItem.Name),
		zap.String("Series ID", jellyfinMovieItem.ID),
	)

	searchResult, err := tmdbAPIClient.SearchMediaByName(
		jellyfinMovieItem.Name,
		int(jellyfinMovieItem.ProductionYear),
		MediaTypeMovie,
	)

	if err != nil {
		// Error is already logged by SearchMediaByName
		enrichMovieWithDefaultInfos(jellyfinMovieItem)
		return
	}

	details := getItemDetailsFromSearchResult(searchResult)
	jellyfinMovieItem.Overview = details.Overview
	jellyfinMovieItem.PosterURL = details.PosterURL
}

func EnrichMovieItemsList(
	jellyfinMovieItem *[]jellyfin.MovieItem,
	tmdbAPIClient APIInterface,
	app *app.ApplicationContext,
) {
	for i := range *jellyfinMovieItem {
		EnrichMovieItem(&(*jellyfinMovieItem)[i], tmdbAPIClient, app)
	}
}
