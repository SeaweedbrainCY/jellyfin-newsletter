package tmdb

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"go.uber.org/zap"
)

func GetMovieDetails(
	tmdbAPIClient APIClient,
	jellyfinMovieItem jellyfin.MovieItem,
	app app.ApplicationContext,
) *ItemDetails {
	if jellyfinMovieItem.TMDBId != "" {
		parsedHTTPResponse, err := tmdbAPIClient.GetMediaByID(jellyfinMovieItem.TMDBId, MediaTypeMovie)

		if err != nil {
			// Error is already logged by GetMediaByID
			return getDefaultItemDetails()
		}

		return getItemDetailsFromHTTPResponse(parsedHTTPResponse)
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
		return getDefaultItemDetails()
	}

	return getItemDetailsFromSearchResult(searchResult)
}
