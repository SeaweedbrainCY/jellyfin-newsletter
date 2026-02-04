package jellyfin

import (
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

func testSeriesInitApp() (*app.ApplicationContext, *observer.ObservedLogs) {
	observedDays := 30
	loggerCore, recordedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(loggerCore)
	return &app.ApplicationContext{
		Logger: logger,
		Config: &config.Configuration{
			Jellyfin: config.JellyfinConfig{
				ObservedPeriodDays:   observedDays,
				WatchedSeriesFolders: []string{"folder1"},
			},
		},
	}, recordedLogs
}

func getSeriesBaseItems() ([]jellyfinAPI.BaseItemDto, []jellyfinAPI.BaseItemDto, []jellyfinAPI.BaseItemDto, []jellyfinAPI.BaseItemDto) {
	recentSeries := []jellyfinAPI.BaseItemDto{
		{
			Id:             Ptr("1813f4b17e9d4a799641c09319b5ffcc"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Series 1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -7))),
			ProviderIds:    map[string]string{"Tmdb": "1027", "Imdb": "2276"},
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		},
		{
			Id:             Ptr("05bc80140f6c41d1a6366c81dac444a4"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Series 2")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -10))),
			ProviderIds:    map[string]string{"Tmdb": "1027", "Imdb": "2276"},
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
		},
		{
			Id:             Ptr("f4971e32089041f3a3d6774277c2ccb9"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -7))),
			ProviderIds:    map[string]string{"Tmdb": "1027", "Imdb": "2276"},
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("1813f4b17e9d4a799641c09319b5ffcc")),
		},
		{
			Id:             Ptr("bcedb6a404974245b41fe224f31e6460"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -7))),
			ProviderIds:    map[string]string{"Tmdb": "1027", "Imdb": "2276"},
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("1813f4b17e9d4a799641c09319b5ffcc")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("f4971e32089041f3a3d6774277c2ccb9")),
		},
		{
			Id:             Ptr("aa1111"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Series 3")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2024))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -5))),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		},
		{
			Id:           Ptr("aa1111-s1"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -4))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("aa1111")),
		},
		{
			Id:           Ptr("aa1111-s1-e1"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Episode 1")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -3))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("aa1111")),
			SeasonId:     *jellyfinAPI.NewNullableString(Ptr("aa1111-s1")),
		},
	}

	recentSeasons := []jellyfinAPI.BaseItemDto{
		// Old series
		{
			Id:           Ptr("bb2222"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Old Series 1")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -90))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		},

		// Recent season
		{
			Id:           Ptr("bb2222-s2"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Season 2")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -30))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("bb2222")),
		},

		// Episode in that season
		{
			Id:           Ptr("bb2222-s2-e1"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Episode 1")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -8))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("bb2222")),
			SeasonId:     *jellyfinAPI.NewNullableString(Ptr("bb2222-s2")),
		},
	}

	recentEpisodes := []jellyfinAPI.BaseItemDto{
		// Old series
		{
			Id:           Ptr("cc3333"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Very Old Series")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -180))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		},

		// Old seasons
		{
			Id:           Ptr("cc3333-s1"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -150))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("cc3333")),
		},
		{
			Id:           Ptr("cc3333-s2"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Season 2")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -140))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("cc3333")),
		},

		// Recent episode
		{
			Id:           Ptr("cc3333-s1-e5"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Episode 5")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -2))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("cc3333")),
			SeasonId:     *jellyfinAPI.NewNullableString(Ptr("cc3333-s1")),
		},
		{
			Id:           Ptr("cc3333-s2-e5"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Episode 5")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -1))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("cc3333")),
			SeasonId:     *jellyfinAPI.NewNullableString(Ptr("cc3333-s2")),
		},
	}

	olderLibrary := []jellyfinAPI.BaseItemDto{
		// Old series
		{
			Id:             Ptr("dd4444"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Legacy Series")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2019))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -120))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		},

		// Old season
		{
			Id:           Ptr("dd4444-s1"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -110))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("dd4444")),
		},

		// Old episode
		{
			Id:           Ptr("dd4444-s1-e1"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Episode 1")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -100))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("dd4444")),
			SeasonId:     *jellyfinAPI.NewNullableString(Ptr("dd4444-s1")),
		},

		// Another old episode (for volume / ordering tests)
		{
			Id:           Ptr("dd4444-s1-e2"),
			Name:         *jellyfinAPI.NewNullableString(Ptr("Episode 2")),
			DateCreated:  *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -95))),
			Type:         Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType: *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:     *jellyfinAPI.NewNullableString(Ptr("dd4444")),
			SeasonId:     *jellyfinAPI.NewNullableString(Ptr("dd4444-s1")),
		},
	}

	return recentSeries, recentSeasons, recentEpisodes, olderLibrary
}

func TestGetNewlyAddedSeries(t *testing.T) {
	mockedApp, recordedLogs := testSeriesInitApp()
	recentSeries, recentSeasons, recentEpisodes, olderLibrary := getSeriesBaseItems()

	allSeriesItems := make(
		[]jellyfinAPI.BaseItemDto,
		len(recentSeries),
		len(recentSeries)+len(recentSeasons)+len(recentEpisodes)+len(olderLibrary),
	)
	_ = copy(allSeriesItems, recentSeries)
	allSeriesItems = append(allSeriesItems, recentSeasons...)
	allSeriesItems = append(allSeriesItems, recentEpisodes...)
	allSeriesItems = append(allSeriesItems, olderLibrary...)

	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetAllItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			return &allSeriesItems, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}

	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	newSeriesItems, err := client.GetNewlyAddedSeries(mockedApp)

	require.NoError(t, err)
	assert.Empty(t, recordedLogs)
	listOfMissingItemsId := []string{}
	for _,series := range recentSeries {
		
	}
}
