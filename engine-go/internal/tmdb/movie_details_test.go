package tmdb

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

// mockHTTPClient lets us simulate HTTP-level failures (e.g. connection refused)
type mockHTTPClient struct {
	err  error
	resp *http.Response
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

// errReader simulates a body that fails on Read
type errReader struct{}

func (e errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("simulated read error")
}

func (e errReader) Close() error { return nil }

func newObservedClient(level zap.AtomicLevel) (APIClient, *observer.ObservedLogs) {
	loggerCore, logs := observer.New(level.Level())
	logger := zap.New(loggerCore)
	return APIClient{
		APIKey:     "test_token",
		Lang:       "en",
		Logger:     logger,
		BaseURL:    "https://api.themoviedb.org/3",
		HTTPClient: http.DefaultClient, // overridden per test
	}, logs
}

// ── GetMediaByID ─────────────────────────────────────────────────────────────

func TestGetMediaByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mediaType     MediaType
		setupServer   func() (*httptest.Server, *mockHTTPClient) // return one or the other
		useRealServer bool
		wantErr       bool
		wantOverview  string
		wantLogCount  int
	}{
		{
			name:      "success - valid response",
			id:        "1234",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"overview":"Great movie","poster_path":"/p.jpg","popularity":9.5}`))
				}))
				return srv, nil
			},
			wantErr:      false,
			wantOverview: "Great movie",
		},
		{
			name:      "HTTP transport error",
			id:        "1234",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				// Return a mock that simulates a connection failure.
				// We also need a dummy server so checkHTTPResponse doesn't panic on nil resp.
				// Instead we use the mock path — server is nil, client handles it.
				mock := &mockHTTPClient{
					err: errors.New("connection refused"),
					// checkHTTPResponse dereferences resp even when err != nil,
					// so we provide a minimal non-nil response to avoid a panic.
					resp: &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{}`)),
					},
				}
				return nil, mock
			},
			wantErr: true,
		},
		{
			name:      "non-2xx status code - 404",
			id:        "9999",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"status_message":"The resource you requested could not be found."}`))
				}))
				return srv, nil
			},
			wantErr: true,
		},
		{
			name:      "non-2xx status code - 401 unauthorized",
			id:        "1234",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"status_message":"Invalid API key"}`))
				}))
				return srv, nil
			},
			wantErr: true,
		},
		{
			name:      "non-2xx status code - 500 server error",
			id:        "1234",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				return srv, nil
			},
			wantErr: true,
		},
		{
			name:      "body read error",
			id:        "1234",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				mock := &mockHTTPClient{
					resp: &http.Response{
						StatusCode: http.StatusOK,
						Body:       errReader{},
					},
				}
				return nil, mock
			},
			// json.Unmarshal on empty bytes returns an error
			wantErr: true,
		},
		{
			name:      "invalid JSON response",
			id:        "1234",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`not valid json {{`))
				}))
				return srv, nil
			},
			wantErr: true,
		},
		{
			name:      "empty JSON object - no fields populated",
			id:        "1234",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{}`))
				}))
				return srv, nil
			},
			wantErr:      false,
			wantOverview: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loggerCore, _ := observer.New(zap.DebugLevel)
			logger := zap.New(loggerCore)

			srv, mock := tc.setupServer()

			var client APIClient
			if srv != nil {
				defer srv.Close()
				client = getTestClient(logger, srv)
			} else {
				client = APIClient{
					APIKey:     "test_token",
					Lang:       "en",
					Logger:     logger,
					BaseURL:    "https://api.themoviedb.org/3",
					HTTPClient: mock,
				}
			}

			result, err := client.GetMediaByID(tc.id, tc.mediaType)

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.wantOverview, result.Overview)
			}
		})
	}
}

// ── SearchMediaByName ─────────────────────────────────────────────────────────

func TestSearchMediaByName(t *testing.T) {
	tests := []struct {
		name           string
		queryName      string
		productionYear int
		mediaType      MediaType
		setupServer    func() (*httptest.Server, *mockHTTPClient)
		wantErr        bool
		wantResultsLen int
	}{
		{
			name:           "success - results returned",
			queryName:      "Inception",
			productionYear: 2010,
			mediaType:      MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write(
						[]byte(
							`{"results":[{"overview":"A mind-bending thriller","popularity":8.2,"poster_path":"/inc.jpg"}]}`,
						),
					)
				}))
				return srv, nil
			},
			wantErr:        false,
			wantResultsLen: 1,
		},
		{
			name:           "success - year omitted when zero",
			queryName:      "Inception",
			productionYear: 0,
			mediaType:      MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Assert the year query param is absent
					assert.Empty(t, r.URL.Query().Get("year"))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"results":[]}`))
				}))
				return srv, nil
			},
			wantErr:        false,
			wantResultsLen: 0,
		},
		{
			name:           "empty name - returns error immediately",
			queryName:      "",
			productionYear: 2020,
			mediaType:      MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				// Server should never be called
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Error("server should not have been called for an empty name")
				}))
				return srv, nil
			},
			wantErr: true,
		},
		{
			name:      "HTTP transport error",
			queryName: "Inception",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				mock := &mockHTTPClient{
					err: errors.New("connection refused"),
					resp: &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{}`)),
					},
				}
				return nil, mock
			},
			wantErr: true,
		},
		{
			name:      "non-2xx status code - 403",
			queryName: "Inception",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"status_message":"Access denied"}`))
				}))
				return srv, nil
			},
			wantErr: true,
		},
		{
			name:      "body read error",
			queryName: "Inception",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				mock := &mockHTTPClient{
					resp: &http.Response{
						StatusCode: http.StatusOK,
						Body:       errReader{},
					},
				}
				return nil, mock
			},
			wantErr: true,
		},
		{
			name:      "invalid JSON response",
			queryName: "Inception",
			mediaType: MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{bad json`))
				}))
				return srv, nil
			},
			wantErr: true,
		},
		{
			name:           "empty results array",
			queryName:      "UnknownMovie12345",
			productionYear: 1900,
			mediaType:      MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"results":[]}`))
				}))
				return srv, nil
			},
			wantErr:        false,
			wantResultsLen: 0,
		},
		{
			name:           "multiple results returned",
			queryName:      "Batman",
			productionYear: 0,
			mediaType:      MediaTypeMovie,
			setupServer: func() (*httptest.Server, *mockHTTPClient) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"results":[
						{"overview":"Batman 1","popularity":5.0,"poster_path":"/b1.jpg"},
						{"overview":"Batman 2","popularity":9.0,"poster_path":"/b2.jpg"}
					]}`))
				}))
				return srv, nil
			},
			wantErr:        false,
			wantResultsLen: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loggerCore, _ := observer.New(zap.DebugLevel)
			logger := zap.New(loggerCore)

			srv, mock := tc.setupServer()

			var client APIClient
			if srv != nil {
				defer srv.Close()
				client = getTestClient(logger, srv)
			} else {
				client = APIClient{
					APIKey:     "test_token",
					Lang:       "en",
					Logger:     logger,
					BaseURL:    "https://api.themoviedb.org/3",
					HTTPClient: mock,
				}
			}

			result, err := client.SearchMediaByName(tc.queryName, tc.productionYear, tc.mediaType)

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, result.Results, tc.wantResultsLen)
			}
		})
	}
}
