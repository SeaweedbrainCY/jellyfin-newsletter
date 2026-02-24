package tmdb

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"go.uber.org/zap"
)

type APIInterface interface {
	GetMediaByID(id string, mediaType MediaType) (*GetMediaHTTPResponse, error)
	SearchMediaByName(name string, productionYear int, mediaType MediaType) (*SearchMediaHTTPResponse, error)
}

type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	APIKey     string
	Lang       string
	Logger     *zap.Logger
	BaseURL    string
	HTTPClient HTTPClientInterface
}

func InitTMDBApiClient(app *app.ApplicationContext) APIClient {
	return APIClient{
		APIKey:     app.Config.TMDB.APIKey,
		Lang:       app.Config.EmailTemplate.Language,
		Logger:     app.Logger,
		BaseURL:    "https://api.themoviedb.org/3",
		HTTPClient: http.DefaultClient,
	}
}

type GetMediaHTTPResponse struct {
	Overview   string  `json:"overview"`
	PosterPath string  `json:"poster_path"`
	Popularity float64 `json:"popularity"`
}

type SearchMediaHTTPResponse struct {
	Results []GetMediaHTTPResponse `json:"results"`
}

func (client APIClient) prepareGetAPIRequest(url string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		client.Logger.Error(
			"An error occurred while building the request towards the TMDB API.",
			zap.Error(err),
		)
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+client.APIKey)

	return request, nil
}

func checkHTTPResponse(url string, resp *http.Response, httpErr error, logger *zap.Logger) error {
	if httpErr != nil {
		logger.Error(
			"HTTP request failed",
			zap.String("URL", url),
			zap.Error(httpErr),
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, bodyErr := io.ReadAll(resp.Body)

		fields := []zap.Field{
			zap.Int("status", resp.StatusCode),
			zap.String("url", url),
		}

		if bodyErr == nil {
			fields = append(fields, zap.String("body", string(body)))
		} else {
			fields = append(fields, zap.NamedError("body_read_error", bodyErr))
		}

		logger.Error("Unexpected HTTP response", fields...)
		return errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

func (client APIClient) GetMediaByID(id string, mediaType MediaType) (*GetMediaHTTPResponse, error) {
	url := client.BaseURL + "/" + mediaType.ToString() + "/" + id + "?language=" + client.Lang

	request, err := client.prepareGetAPIRequest(url)

	if err != nil {
		return nil, err
	}

	httpResponse, execReqErr := client.HTTPClient.Do(request)

	err = checkHTTPResponse(url, httpResponse, execReqErr, client.Logger)

	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	body, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		client.Logger.Error("Impossible to read the HTTP response body.",
			zap.String("URL", url),
			zap.Int("HTTP Status code", httpResponse.StatusCode),
			zap.Error(err))
	}

	var decodedBody GetMediaHTTPResponse
	jsonDecodeErr := json.Unmarshal(body, &decodedBody)

	if jsonDecodeErr != nil {
		client.Logger.Error(
			"An error occurred while decoding TMDB API's answer.",
			zap.Error(jsonDecodeErr),
			zap.String("URL", url),
		)
		return nil, jsonDecodeErr
	}
	return &decodedBody, nil
}

func (client APIClient) SearchMediaByName(
	name string,
	productionYear int,
	mediaType MediaType,
) (*SearchMediaHTTPResponse, error) {
	if name == "" {
		client.Logger.Warn(
			"Attempted to search for an item on TMDB but the given item had an Unknown Name. Operation has been aborted",
		)
		return nil, errors.New("empty name")
	}

	url := client.BaseURL + "/search/" + mediaType.ToString() + "?query=" + name + "&language=" + client.Lang

	if productionYear != 0 {
		url += "&year=" + strconv.Itoa(productionYear)
	}

	request, err := client.prepareGetAPIRequest(url)

	if err != nil {
		return nil, err
	}

	httpResponse, execReqErr := client.HTTPClient.Do(request)

	err = checkHTTPResponse(url, httpResponse, execReqErr, client.Logger)

	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	body, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		client.Logger.Error("Impossible to read the HTTP response body.",
			zap.String("URL", url),
			zap.Int("HTTP Status code", httpResponse.StatusCode),
			zap.Error(err))
	}

	var decodedBody SearchMediaHTTPResponse
	jsonDecodeErr := json.Unmarshal(body, &decodedBody)

	if jsonDecodeErr != nil {
		client.Logger.Error(
			"An error occurred while decoding TMDB API's answer.",
			zap.Error(jsonDecodeErr),
			zap.String("URL", url),
		)
		return nil, jsonDecodeErr
	}
	return &decodedBody, nil
}
