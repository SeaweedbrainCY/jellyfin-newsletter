package tmdb

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"go.uber.org/zap"
)

func GetSeriesDetails(
	tmdbAPIClient APIClient,
	jellyfinSeriesItem jellyfin.NewlyAddedSeriesItem,
	app app.ApplicationContext,
) *ItemDetails {
	if jellyfinSeriesItem.TMDBId != "" {
		parsedHTTPResponse, err := tmdbAPIClient.GetMediaByID(jellyfinSeriesItem.TMDBId, MediaTypeSeries)

		if err != nil {
			// Error is already logged by GetMediaByID
			return getDefaultItemDetails()
		}

		return getItemDetailsFromHTTPResponse(parsedHTTPResponse)
	}
	// No TMDB id, we perform a search by name and select the item with the highest popularity
	app.Logger.Debug(
		"Series has no TMDB id. TMDB information will be retrieved by searching with Series's name. If several media match, the choice will be based on the highest popularity.",
		zap.String("Series Name", jellyfinSeriesItem.SeriesName),
		zap.String("Series ID", jellyfinSeriesItem.SeriesID),
	)

	searchResult, err := tmdbAPIClient.SearchMediaByName(
		jellyfinSeriesItem.SeriesName,
		jellyfinSeriesItem.ProductionYear,
		MediaTypeSeries,
	)

	if err != nil {
		// Error is already logged by SearchMediaByName
		return getDefaultItemDetails()
	}

	return getItemDetailsFromSearchResult(searchResult)

}
