package jellyfin

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type MockJellyfinItemsAPI struct {
	ExecuteGetMoviesItemsByFolderID func() (*[]jellyfinAPI.BaseItemDto, error)
	ExecuteGetRootFolderIDByName    func() (string, error)
}

func (m MockJellyfinItemsAPI) GetMoviesItemsByFolderID(
	_ string,
	_ bool,
	_ *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return m.ExecuteGetMoviesItemsByFolderID()
}

func (m MockJellyfinItemsAPI) GetAllItemsByFolderID(
	_ string,
	_ *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return &[]jellyfinAPI.BaseItemDto{}, nil
}

func (m MockJellyfinItemsAPI) GetRootFolderIDByName(_ string, _ *app.ApplicationContext) (string, error) {
	return m.ExecuteGetRootFolderIDByName()
}

func Ptr[T any](v T) *T {
	return &v
}

func initApp() (*app.ApplicationContext, *observer.ObservedLogs) {
	observedDays := 30
	loggerCore, recordedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(loggerCore)
	return &app.ApplicationContext{
		Logger: logger,
		Config: &config.Configuration{
			Jellyfin: config.JellyfinConfig{
				ObservedPeriodDays: observedDays,
			},
		},
	}, recordedLogs
}

func TestGetRecentlyAddedMoviesByFolder(t *testing.T) {
	mockedApp, recordedLogs := initApp()
	mockedRecentlyAddedItems := []jellyfinAPI.BaseItemDto{
		{
			Id:             Ptr("8b54388aca994d4fb867944d3150a7e0"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2026))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now())),
			ProviderIds:    map[string]string{"Tmdb": "2876", "Imdb": "2653"},
		},
		{
			Id:             Ptr("2264338b9fe1475b8f2b8095531dd1ff"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 2")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -7))),
			ProviderIds:    map[string]string{"Tmdb": "1027", "Imdb": "2276"},
		},
		{
			Id:             Ptr("7c48cba998264cea90880b9efd0d9c9b"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 3")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2022))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -30.))),
			ProviderIds:    map[string]string{"Tmdb": "1092"},
		},
	}

	mockedOtherItems := []jellyfinAPI.BaseItemDto{
		{
			Id:             Ptr("613a4ec7d1654206b8bc2130b1e785ce"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 4")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2026))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -31))),
			ProviderIds:    map[string]string{"Tmdb": "2173"},
		},
		{
			Id:             Ptr("ee1b42664d9a41ca9d7db8ca114f3fdb"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 5")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, -2, 0))),
			ProviderIds:    map[string]string{"Tmdb": "9876"},
		},
		{
			Id:             Ptr("bd9fb02655f54c79932431a342763c8b"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 6")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2021))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(-1, 0, 0))),
			ProviderIds:    map[string]string{"Tmdb": "1098"},
		},
	}

	concatMovies := make(
		[]jellyfinAPI.BaseItemDto,
		len(mockedRecentlyAddedItems),
		len(mockedRecentlyAddedItems)+len(mockedOtherItems),
	)
	_ = copy(concatMovies, mockedRecentlyAddedItems)
	concatMovies = append(concatMovies, mockedOtherItems...)
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			return &concatMovies, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newlyAddedMovies, err := client.getRecentlyAddedMoviesByFolder("folderName", mockedApp)

	require.NoError(t, err)
	assert.Empty(t, recordedLogs)
	assert.Len(
		t,
		mockedRecentlyAddedItems,
		len(newlyAddedMovies),
		"Expected : %v. Actual: %v",
		mockedRecentlyAddedItems,
		newlyAddedMovies,
	)
	for _, expectedMovie := range mockedRecentlyAddedItems {
		isMoviePresent := false
		for _, movie := range newlyAddedMovies {
			if movie.ID == *expectedMovie.Id {
				isMoviePresent = true
				assert.Equal(
					t,
					*expectedMovie.Name.Get(),
					movie.Name,
					"Movie ID %s: Not the right movie name. Expected: %s. Actual: %s",
					movie.ID,
					expectedMovie.Name.Get(),
					movie.Name,
				)
				assert.Equal(
					t,
					*expectedMovie.ProductionYear.Get(),
					movie.ProductionYear,
					"Movie ID %s: Not the right production year. Expected: %d. Actual: %d",
					movie.ID,
					expectedMovie.ProductionYear.Get(),
					movie.ProductionYear,
				)
				assert.Equal(
					t,
					expectedMovie.DateCreated.Get(),
					movie.AdditionDate,
					"Movie ID %s: Not the right creation date. Expected: %s. Actual: %s",
					movie.ID,
					expectedMovie.DateCreated.Get().String(),
					movie.AdditionDate.String(),
				)
				expectedProviderID, atoiErr := strconv.Atoi(expectedMovie.ProviderIds["Tmdb"])
				require.NoError(t, atoiErr)
				assert.Equal(
					t,
					expectedProviderID,
					movie.TMDBId,
					"Movie ID %s: Not the TMDBID. Expected: %s. Actual: %d",
					movie.ID,
					expectedMovie.ProviderIds["Tmdb"],
					movie.TMDBId,
				)
				break
			}
		}
		assert.True(
			t,
			isMoviePresent,
			"Movie %s is not present in the actual recently added movie.",
			*expectedMovie.Id,
		)
	}
}

func baseMovie() []jellyfinAPI.BaseItemDto {
	return []jellyfinAPI.BaseItemDto{
		{
			Id:             Ptr("8b54388aca994d4fb867944d3150a7e0"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2026))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now())),
			ProviderIds:    map[string]string{"Tmdb": "2876", "Imdb": "2653"},
		},
		{
			Id:             Ptr("2264338b9fe1475b8f2b8095531dd1ff"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 2")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -7))),
			ProviderIds:    map[string]string{"Tmdb": "1027", "Imdb": "2276"},
		},
		{
			Id:             Ptr("7c48cba998264cea90880b9efd0d9c9b"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Movie 3")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2022))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, -2, 0.))),
			ProviderIds:    map[string]string{"Tmdb": "1092"},
		},
	}
}

func TestGetRecentlyAddedMoviesByFolderWithMovieNameNull(t *testing.T) {
	app, recordedLogs := initApp()
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			items := baseMovie()
			items[0].Name = *jellyfinAPI.NewNullableString(nil)
			return &items, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newlyAddedMovies, err := client.getRecentlyAddedMoviesByFolder("folderName", app)

	require.NoError(t, err)
	assert.Empty(t, recordedLogs)
	assert.Len(t, newlyAddedMovies, 2)
	isMovieNameDefault := false
	for _, movie := range newlyAddedMovies {
		if movie.ID == "8b54388aca994d4fb867944d3150a7e0" {
			assert.Equal(t, "Unknown Movie Name", movie.Name)
			isMovieNameDefault = true
		}
	}
	assert.True(t, isMovieNameDefault, "No movie with ID 8b54388aca994d4fb867944d3150a7e0 found. This is not expected.")
}

func TestGetRecentlyAddedMoviesByFolderWithMovieProductionYearNull(t *testing.T) {
	app, recordedLogs := initApp()
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			items := baseMovie()
			items[0].ProductionYear = *jellyfinAPI.NewNullableInt32(nil)
			return &items, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newlyAddedMovies, err := client.getRecentlyAddedMoviesByFolder("folderName", app)

	require.NoError(t, err)
	assert.Equal(t, 0, recordedLogs.Len())
	assert.Len(t, newlyAddedMovies, 2)
	isMovieNameDefault := false
	for _, movie := range newlyAddedMovies {
		if movie.ID == "8b54388aca994d4fb867944d3150a7e0" {
			assert.Equal(t, int32(0), movie.ProductionYear)
			isMovieNameDefault = true
		}
	}
	assert.True(t, isMovieNameDefault, "No movie with ID 8b54388aca994d4fb867944d3150a7e0 found. This is not expected.")
}

func TestGetRecentlyAddedMoviesByFolderWithNoTMDBID(t *testing.T) {
	app, recordedLogs := initApp()
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			items := baseMovie()
			items[0].ProviderIds = map[string]string{}
			return &items, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newlyAddedMovies, err := client.getRecentlyAddedMoviesByFolder("folderName", app)

	require.NoError(t, err)
	assert.Equal(t, 0, recordedLogs.Len())
	assert.Len(t, newlyAddedMovies, 2)
	isMovieNameDefault := false
	for _, movie := range newlyAddedMovies {
		if movie.ID == "8b54388aca994d4fb867944d3150a7e0" {
			assert.Equal(t, 0, movie.TMDBId)
			isMovieNameDefault = true
		}
	}
	assert.True(t, isMovieNameDefault, "No movie with ID 8b54388aca994d4fb867944d3150a7e0 found. This is not expected.")
}

func TestGetRecentlyAddedMoviesByFolderWithNoCreationDate(t *testing.T) {
	app, recordedLogs := initApp()
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			items := baseMovie()
			items[0].DateCreated = *jellyfinAPI.NewNullableTime(nil)
			return &items, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newlyAddedMovies, err := client.getRecentlyAddedMoviesByFolder("folderName", app)

	require.NoError(t, err)

	require.Equal(t, 1, recordedLogs.Len())
	assert.Equal(t, "Movie ignored because it has no creation date.", recordedLogs.All()[0].Message)

	fields := recordedLogs.All()[0].Context
	require.Len(t, fields, 2)
	assert.Equal(t, "MovieID", fields[0].Key)
	assert.Equal(t, "8b54388aca994d4fb867944d3150a7e0", fields[0].String)

	assert.Len(t, newlyAddedMovies, 1)
}

func TestGetRecentlyAddedMoviesByFolderWithNoMovies(t *testing.T) {
	app, recordedLogs := initApp()
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			return &[]jellyfinAPI.BaseItemDto{}, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newlyAddedMovies, err := client.getRecentlyAddedMoviesByFolder("folderName", app)

	require.NoError(t, err)

	require.Equal(t, 0, recordedLogs.Len())
	assert.Empty(t, newlyAddedMovies)
}

func TestGetRecentlyAddedMoviesByFolderWithErrorWhileRetrievingFolder(t *testing.T) {
	app, _ := initApp()
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			items := baseMovie()
			return &items, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "", errors.New("error")
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	_, err := client.getRecentlyAddedMoviesByFolder("folderName", app)

	require.Error(t, err)
}

func TestGetRecentlyAddedMoviesByFolderWithErrorWhileRetrievingMovies(t *testing.T) {
	app, _ := initApp()
	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			return nil, errors.New("error")
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	_, err := client.getRecentlyAddedMoviesByFolder("folderName", app)

	require.Error(t, err)
}
