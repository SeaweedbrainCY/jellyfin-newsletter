package tmdb

import "github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"

func GetSeriesDetails(
	tmdbAPIClient TMDBAPIClient,
	jellyfinSeriesItem jellyfin.NewlyAddedSeriesItem,
) (*ItemDetails, error) {
	if jellyfinSeriesItem.TMDBId != "" {
		parsedHttpResponse, err := tmdbAPIClient.GetMediaByID(jellyfinSeriesItem.TMDBId, MediaTypeSeries)

		if err != nil {
			// Error is already logged by GetMediaByID
			return nil, err
		}
		return getItemDetailsFromHTTPResponse(parsedHttpResponse), nil
	}
}
