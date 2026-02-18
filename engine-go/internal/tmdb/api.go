package tmdb

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type GetMediaHTTPResponse struct {
	Overview   string `json:"overview"`
	PosterPath string `json:"poster_path"`
}

func (client TMDBAPIClient) executeGetMediaRequest(url string) (*GetMediaHTTPResponse, error) {
	request, buildReqErr := http.NewRequest(http.MethodGet, url, nil)

	if buildReqErr != nil {
		client.Logger.Error(
			"An error occurred while building the request towards the TMDB API.",
			zap.Error(buildReqErr),
		)
		return nil, buildReqErr
	}

	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+client.APIKey)

	httpResponse, execReqErr := http.DefaultClient.Do(request)

	body, bodyErr := io.ReadAll(httpResponse.Body)

	if execReqErr != nil || bodyErr != nil {
		bodyString := ""
		if bodyErr == nil {
			bodyString = string(body)
			defer httpResponse.Body.Close()
		}
		client.Logger.Error(
			"An error occurred while requesting the TMDB API.",
			zap.Int("HTTP status", httpResponse.StatusCode),
			zap.String("HTTP response body", bodyString),
			zap.Bool("Is body readable", bodyErr != nil),
			zap.Error(execReqErr),
		)
		return nil, execReqErr
	}
	defer httpResponse.Body.Close()

	var decodedBody GetMediaHTTPResponse
	jsonDecodeErr := json.Unmarshal(body, &decodedBody)

	if jsonDecodeErr != nil {
		client.Logger.Error(
			"An error occurred while decoding TMDB API's answer.",
			zap.Error(jsonDecodeErr),
		)
		return nil, jsonDecodeErr
	}
	return &decodedBody, nil
}

func (client TMDBAPIClient) GetMediaByID(id string, mediaType MediaType) (*GetMediaHTTPResponse, error) {
	url := "https://api.themoviedb.org/3/" + mediaType.ToString() + "/" + id + "?language=" + client.Lang
	return client.executeGetMediaRequest(url)
}
