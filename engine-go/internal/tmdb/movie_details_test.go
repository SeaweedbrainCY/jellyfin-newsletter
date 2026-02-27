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
	defaultOverview := "No description available."
	defaultPosterURL := "https://placehold.co/200"
	tests := []struct {
		name               string
		tmdbID             string
		testServerHandler  http.Handler
		expectedOverview   string
		expectedPosterPath string
		expectErr          bool
	}{
		{
			name:   "Success - Valid response",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"overview": "This is the description of a media", "popularity": 12.5, "poster_path":"/poster/path"}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:   "Success - No results",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{}`,
					),
				)
			}),
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
			expectErr:          false,
		},
		{
			name:   "Error - Malformed json",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{{"overview": "This is the description of`,
					),
				)
			}),
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
			expectErr:          true,
		},
		{
			name:   "Success - Missing overview",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"popularity": 12.5, "poster_path":"/poster/path"}`,
					),
				)
			}),
			expectedOverview:   defaultOverview,
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:   "Success - Missing posterPath",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"overview": "This is the description of a media", "popularity": 12.5}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: defaultPosterURL,
			expectErr:          false,
		},
		{
			name:   "Success - Missing popularity",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"overview": "This is the description of a media",  "poster_path":"/poster/path"}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:   "Error - Connection reset",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				hj, ok := w.(http.Hijacker)
				if !ok {
					panic("not a hijacker")
				}
				conn, _, _ := hj.Hijack()
				conn.Close()
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:   "Error - partial response EOF",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Content-Length", "1000") // lie about body size
				w.Write([]byte(`{"overview": `))
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:   "Error - Error 404",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:   "Error - Error 403",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:   "Error - Error 500",
			tmdbID: "12345",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testServer := httptest.NewServer(testCase.testServerHandler)
			defer testServer.Close()
			loggerCore, recordedLogs := observer.New(zap.InfoLevel)
			logger := zap.New(loggerCore)
			client := getTestClient(logger, testServer)
			jellyfinMovieItem := getBaseJellyfinMovieItem()
			jellyfinMovieItem.TMDBId = testCase.tmdbID
			app := app.ApplicationContext{
				Logger: logger,
			}
			movieDetails := GetMovieDetails(client, jellyfinMovieItem, app)
			if testCase.expectErr {
				assert.NotEmpty(t, recordedLogs.All())
			} else {
				require.Empty(t, recordedLogs.All())
			}
			assert.Equal(t, testCase.expectedOverview, movieDetails.Overview)
			assert.Equal(t, testCase.expectedPosterPath, movieDetails.PosterURL)
		})
	}
}

func TestGetMovieDetailsWithSearchByName(t *testing.T) {
	defaultOverview := "No description available."
	defaultPosterURL := "https://placehold.co/200"
	tests := []struct {
		name               string
		movieName          string
		testServerHandler  http.Handler
		expectedOverview   string
		expectedPosterPath string
		expectErr          bool
	}{
		{
			name:      "Success - Valid response multiple movies",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"overview": "This is the description of a media", "popularity": 2.8, "poster_path":"/poster/path"}, {"overview": "This is the description of the most popular media", "popularity": 12.5, "poster_path":"/poster/popular/path"}]}`,
					),
				)
			}),
			expectedOverview:   "This is the description of the most popular media",
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/popular/path",
			expectErr:          false,
		},
		{
			name:      "Success - Valid response one movie",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"overview": "This is the description of a media", "popularity": 2.8, "poster_path":"/poster/path"}]}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:      "Success - No results",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": []}`,
					),
				)
			}),
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
			expectErr:          false,
		},
		{
			name:      "Error - Malformed json",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"overview": "This is the description of`,
					),
				)
			}),
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
			expectErr:          true,
		},
		{
			name:      "Success - Missing overview",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"popularity": 2.8, "poster_path":"/poster/path"}]}`,
					),
				)
			}),
			expectedOverview:   defaultOverview,
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:      "Success - Missing posterPath",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"overview": "This is the description of a media", "popularity": 2.8}]}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: defaultPosterURL,
			expectErr:          false,
		},
		{
			name:      "Success - Missing popularity",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"overview": "This is the description of a media", "poster_path":"/poster/path"}]}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:      "Success - Missing popularity multiple movies",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"overview": "This is the description of the most popular media", "poster_path":"/poster/popular/path"}, {"overview": "This is the description of a media", "popularity": 2.8, "poster_path":"/poster/path"}]}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:      "Success - Same popularity",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(
					[]byte(
						`{"results": [{"overview": "This is the description of a media", "popularity": 2.8, "poster_path":"/poster/path"}, {"overview": "This is the description of the most popular media","popularity": 2.8, "poster_path":"/poster/popular/path"}]}`,
					),
				)
			}),
			expectedOverview:   "This is the description of a media",
			expectedPosterPath: "https://image.tmdb.org/t/p/w500/poster/path",
			expectErr:          false,
		},
		{
			name:      "Error - Connection reset",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				hj, ok := w.(http.Hijacker)
				if !ok {
					panic("not a hijacker")
				}
				conn, _, _ := hj.Hijack()
				conn.Close()
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:      "Error - partial response EOF",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Content-Length", "1000") // lie about body size
				w.Write([]byte(`{"results":`))
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:      "Error - Error 404",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:      "Error - Error 403",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
		{
			name:      "Error - Error 500",
			movieName: "Test movie",
			testServerHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectErr:          true,
			expectedOverview:   defaultOverview,
			expectedPosterPath: defaultPosterURL,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testServer := httptest.NewServer(testCase.testServerHandler)
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
			if testCase.expectErr {
				assert.NotEmpty(t, recordedLogs.All())
			} else {
				require.Empty(t, recordedLogs.All())
			}
			assert.Equal(t, testCase.expectedOverview, movieDetails.Overview)
			assert.Equal(t, testCase.expectedPosterPath, movieDetails.PosterURL)
		})
	}
}
