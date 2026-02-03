package jellyfin

import (
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
}

func (m MockJellyfinItemsAPI) GetMoviesItemsByFolderID(
	folderID string,
	recursive bool,
	app *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return m.ExecuteGetMoviesItemsByFolderID()
}

func (m MockJellyfinItemsAPI) GetAllItemsByFolderID(
	folderID string,
	app *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return nil, nil
}

func (m MockJellyfinItemsAPI) GetRootFolderIDByName(folderName string, app *app.ApplicationContext) (string, error) {
	return "coucuo", nil
}

func Ptr[T any](v T) *T {
	return &v
}

func TestGetRecentlyAddedMoviesByFolder(t *testing.T) {
	observedDays := 30
	loggerCore, recordedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(loggerCore)
	mockedApp := app.ApplicationContext{
		Logger: logger,
		Config: &config.Configuration{
			Jellyfin: config.JellyfinConfig{
				ObservedPeriodDays: observedDays,
			},
		},
	}

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

	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetMoviesItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			items := append(mockedRecentlyAddedItems, mockedOtherItems...)
			return &items, nil
		},
	}
	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newlyAddedMovies, err := client.getRecentlyAddedMoviesByFolder("folderName", &mockedApp)

	require.NoError(t, err)
	assert.Equal(t, recordedLogs.Len(), 0)
	assert.Equal(
		t,
		len(mockedRecentlyAddedItems),
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
				expectedProviderID, err := strconv.Atoi(expectedMovie.ProviderIds["Tmdb"])
				require.NoError(t, err)
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
