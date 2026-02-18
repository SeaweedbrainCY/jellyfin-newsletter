package tmdb

type ItemDetails struct {
	Overview  string
	PosterURL string
}

func getDefaultItemDetails() *ItemDetails {
	return &ItemDetails{
		Overview:  "No description available.",
		PosterURL: "https://placehold.co/200",
	}
}

func getItemDetailsFromHTTPResponse(parsedHttpResponse *GetMediaHTTPResponse) *ItemDetails {
	itemDetails := getDefaultItemDetails()
	if parsedHttpResponse.Overview != "" {
		itemDetails.Overview = parsedHttpResponse.Overview
	}
	if parsedHttpResponse.PosterPath != "" {
		itemDetails.PosterURL = "https://image.tmdb.org/t/p/w500" + parsedHttpResponse.PosterPath
	}
	return itemDetails
}

// Parse result from the TMDB search based on name
// The final item will be selected based on popularity.
func getItemDetailsFromSearchResult(result *SearchMediaHTTPResponse) *ItemDetails {
	itemDetails := getDefaultItemDetails()
	popularity := 0.0
	for _, item := range result.Results {
		if item.Popularity > popularity {
			if item.Overview != "" {
				itemDetails.Overview = item.Overview
			}
			if item.PosterPath != "" {
				itemDetails.PosterURL = "https://image.tmdb.org/t/p/w500" + item.PosterPath
			}
			popularity = item.Popularity
		}
	}
	return itemDetails
}
