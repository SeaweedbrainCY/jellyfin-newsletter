package tmdb

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func getTestClient(logger *zap.Logger, testServer *httptest.Server) APIClient {
	return APIClient{
		APIKey:     "test_token",
		Lang:       "en",
		Logger:     logger,
		BaseURL:    testServer.URL,
		HTTPClient: testServer.Client(),
	}
}

func getBaseJellyfinMovieItem() jellyfin.MovieItem {
	additionDate := time.Date(2026, 01, 01, 01, 01, 01, 01, time.UTC)
	return jellyfin.MovieItem{
		ID:             "aa1111",
		Name:           "Movie 1",
		AdditionDate:   &additionDate,
		TMDBId:         "1234",
		ProductionYear: 2026,
	}
}
func TestGetMovieDetailsWithTMDBID(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(
			[]byte(
				`{"overview": "This is the description of a media", "popularity": 12.5, "poster_path":"/poster/path"}`,
			),
		)
	}))
	defer testServer.Close()
	loggerCore, recordedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(loggerCore)
	client := getTestClient(logger, testServer)

	jellyfinMovieItem := getBaseJellyfinMovieItem()
	app := app.ApplicationContext{
		Logger: logger,
	}
	movieDetails := GetMovieDetails(client, jellyfinMovieItem, app)

	require.Empty(t, recordedLogs.All())
	assert.Equal(t, "This is the description of a media", movieDetails.Overview)
	assert.Equal(t, "https://image.tmdb.org/t/p/w500/poster/path", movieDetails.PosterURL)
}

func TestGetMovieDetailsWithSearchByName(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(
			[]byte(
				`{"results": [{"overview": "This is the description of a media", "popularity": 2.8, "poster_path":"/poster/path"}, {"overview": "This is the description of the most popular media", "popularity": 12.5, "poster_path":"/poster/popular/path"}]}`,
			),
		)
	}))
	defer testServer.Close()
	loggerCore, recordedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(loggerCore)
	client := getTestClient(logger, testServer)

	jellyfinMovieItem := getBaseJellyfinMovieItem()
	jellyfinMovieItem.TMDBId = ""
	app := app.ApplicationContext{
		Logger: logger,
	}
	movieDetails := GetMovieDetails(client, jellyfinMovieItem, app)

	require.Empty(t, recordedLogs.All())
	assert.Equal(t, "This is the description of the most popular media", movieDetails.Overview)
	assert.Equal(t, "https://image.tmdb.org/t/p/w500/poster/popular/path", movieDetails.PosterURL)
}
