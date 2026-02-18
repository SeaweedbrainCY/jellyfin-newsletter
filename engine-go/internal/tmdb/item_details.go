package tmdb

type ItemDetails struct {
	Overview  string
	PosterURL string
}

func getItemDetailsFromHTTPResponse(parsedHttpResponse *GetMediaHTTPResponse) *ItemDetails {
	itemDetails := ItemDetails{
		Overview:  "No description available.",
		PosterURL: "https://placehold.co/200",
	}
	if parsedHttpResponse.Overview != "" {
		itemDetails.Overview = parsedHttpResponse.Overview
	}
	if parsedHttpResponse.PosterPath != "" {
		itemDetails.PosterURL = "https://image.tmdb.org/t/p/w500" + parsedHttpResponse.PosterPath
	}
	return &itemDetails
}
