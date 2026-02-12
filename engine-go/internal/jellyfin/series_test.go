package jellyfin

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type TableTests struct {
	getSeriesBaseItems            func() []jellyfinAPI.BaseItemDto
	getExpectedResultFromBaseItem func() []NewlyAddedSeriesItem
	name                          string
	loggedMessages                []observer.LoggedEntry
}

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

func getExpectedResultFromBaseItem() []NewlyAddedSeriesItem {
	return []NewlyAddedSeriesItem{
		{
			SeriesName:     "Series 1",
			SeriesID:       "1813f4b17e9d4a799641c09319b5ffcc",
			IsSeriesNew:    true,
			NewSeasons:     nil,
			TMDBId:         1027,
			ProductionYear: 2023,
			AdditionDate:   time.Now().AddDate(0, 0, -7),
		},
		{
			SeriesName:     "Series 2",
			SeriesID:       "aa1111",
			IsSeriesNew:    true,
			NewSeasons:     nil,
			TMDBId:         3001,
			ProductionYear: 2024,
			AdditionDate:   time.Now().AddDate(0, 0, -5),
		},
		{
			SeriesName:  "Old Series 1",
			SeriesID:    "bb2222",
			IsSeriesNew: false,
			NewSeasons: map[string]SeasonItem{
				"bb2222-s2": {
					Name:         "Season 2",
					AdditionDate: time.Now().AddDate(0, 0, -30),
					SeasonNumber: 2,
					Episodes:     nil,
					IsSeasonNew:  true,
				},
			},
			TMDBId:         3001,
			ProductionYear: 2023,
			AdditionDate:   time.Now().AddDate(0, 0, -90),
		},
		{
			SeriesName:  "Very Old Series",
			SeriesID:    "cc3333",
			IsSeriesNew: false,
			NewSeasons: map[string]SeasonItem{
				"cc3333-s1": {
					Name:         "Season 1",
					AdditionDate: time.Now().AddDate(0, 0, -150),
					SeasonNumber: 1,
					IsSeasonNew:  false,
					Episodes: map[string]EpisodeItem{
						"cc3333-s1-e5": {
							Name:          "Episode 5",
							AdditionDate:  time.Now().AddDate(0, 0, -2),
							EpisodeNumber: 5,
						},
					},
				},
				"cc3333-s2": {
					Name:         "Season 2",
					AdditionDate: time.Now().AddDate(0, 0, -140),
					SeasonNumber: 2,
					IsSeasonNew:  false,
					Episodes: map[string]EpisodeItem{
						"cc3333-s2-e5": {
							Name:          "Episode 5",
							AdditionDate:  time.Now().AddDate(0, 0, -1),
							EpisodeNumber: 5,
						},
					},
				},
			},
			TMDBId:         3001,
			ProductionYear: 2023,
			AdditionDate:   time.Now().AddDate(0, 0, -180),
		},
		{
			SeriesName:  "Very Old Series",
			SeriesID:    "ee5555",
			IsSeriesNew: false,
			NewSeasons: map[string]SeasonItem{
				"ee5555-s1": {
					Name:         "Season 1",
					AdditionDate: time.Now().AddDate(0, 0, -150),
					SeasonNumber: 1,
					IsSeasonNew:  false,
					Episodes: map[string]EpisodeItem{
						"ee5555-s1-e5": {
							Name:          "Episode 5",
							AdditionDate:  time.Now().AddDate(0, 0, -2),
							EpisodeNumber: 5,
						},
					},
				},
				"ee5555-s2": {
					Name:         "Season 2",
					AdditionDate: time.Now().AddDate(0, 0, -5),
					SeasonNumber: 2,
					IsSeasonNew:  true,
					Episodes:     nil,
				},
			},
			TMDBId:         3001,
			ProductionYear: 2023,
			AdditionDate:   time.Now().AddDate(0, 0, -180),
		},
	}
}

func getSeriesBaseItems() []jellyfinAPI.BaseItemDto {
	return []jellyfinAPI.BaseItemDto{
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
			Id:             Ptr("f4971e32089041f3a3d6774277c2ccb9"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -7))),
			ProviderIds:    map[string]string{"Tmdb": "1027", "Imdb": "2276"},
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("1813f4b17e9d4a799641c09319b5ffcc")),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
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
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},
		{
			Id:             Ptr("aa1111"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Series 2")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2024))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -5))),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		},
		{
			Id:             Ptr("aa1111-s1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -4))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("aa1111")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},
		{
			Id:             Ptr("aa1111-s1-e1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -3))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("aa1111")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("aa1111-s1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},

		// Old series but recent season
		{
			Id:             Ptr("bb2222"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Old Series 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -90))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
		},
		// old season
		{
			Id:             Ptr("bb2222-s1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -50))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("bb2222")),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
			SeriesName:     *jellyfinAPI.NewNullableString(Ptr("Old Series 1")),
		},

		// Recent season
		{
			Id:             Ptr("bb2222-s2"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 2")),
			SeriesName:     *jellyfinAPI.NewNullableString(Ptr("Old Series 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -30))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("bb2222")),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(2))),
		},

		// Episode in that season
		{
			Id:             Ptr("bb2222-s2-e1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -8))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("bb2222")),
			SeriesName:     *jellyfinAPI.NewNullableString(Ptr("Old Series 1")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("bb2222-s2")),
			SeasonName:     *jellyfinAPI.NewNullableString(Ptr("Season 2")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},

		// Old series but recent episodes
		{
			Id:             Ptr("cc3333"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Very Old Series")),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -180))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
		},

		// Old seasons
		{
			Id:             Ptr("cc3333-s1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -150))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("cc3333")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},
		{
			Id:             Ptr("cc3333-s2"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 2")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -140))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("cc3333")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(2))),
		},

		// Recent episode
		{
			Id:             Ptr("cc3333-s1-e5"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 5")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -2))),
			SeasonName:     *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			SeriesName:     *jellyfinAPI.NewNullableString(Ptr("Very Old Series")),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("cc3333")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("cc3333-s1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(5))),
		},
		{
			Id:             Ptr("cc3333-s2-e5"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 5")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -1))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("cc3333")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("cc3333-s2")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(5))),
		},
		// series with recent episodes and recent seasons
		{
			Id:             Ptr("ee5555"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Very Old Series")),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -180))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
		},

		// Old seasons
		{
			Id:             Ptr("ee5555-s1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -150))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("ee5555")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},
		{
			Id:             Ptr("ee5555-s2"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 2")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -5))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("ee5555")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(2))),
		},

		{
			Id:             Ptr("ee5555-s1-e5"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 5")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -2))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("ee5555")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("ee5555-s1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(5))),
		},
		{
			Id:             Ptr("ee5555-s2-e5"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 5")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -1))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("ee5555")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("ee5555-s2")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(5))),
		},
		// Old series
		{
			Id:             Ptr("dd4444"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Legacy Series")),
			ProviderIds:    map[string]string{"Tmdb": "3001"},
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2019))),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -120))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		},

		// Old season
		{
			Id:             Ptr("dd4444-s1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Season 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -110))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_SEASON),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("dd4444")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},

		// Old episode
		{
			Id:             Ptr("dd4444-s1-e1"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 1")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -100))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("dd4444")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("dd4444-s1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(1))),
		},

		// Another old episode (for volume / ordering tests)
		{
			Id:             Ptr("dd4444-s1-e2"),
			Name:           *jellyfinAPI.NewNullableString(Ptr("Episode 2")),
			DateCreated:    *jellyfinAPI.NewNullableTime(Ptr(time.Now().AddDate(0, 0, -95))),
			Type:           Ptr(jellyfinAPI.BASEITEMKIND_EPISODE),
			LocationType:   *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
			SeriesId:       *jellyfinAPI.NewNullableString(Ptr("dd4444")),
			SeasonId:       *jellyfinAPI.NewNullableString(Ptr("dd4444-s1")),
			ProductionYear: *jellyfinAPI.NewNullableInt32(Ptr(int32(2023))),
			IndexNumber:    *jellyfinAPI.NewNullableInt32(Ptr(int32(2))),
		},
	}
}

func testReturnedSeriesIsCorrect(t *testing.T, expected *NewlyAddedSeriesItem, returned *NewlyAddedSeriesItem) {
	assert.Equal(t, expected.SeriesName, returned.SeriesName, "Series ID %s", expected.SeriesID)
	assert.InDelta(t, expected.AdditionDate.Unix(), returned.AdditionDate.Unix(), 10, "Series ID %s", expected.SeriesID)
	require.Equal(t, expected.IsSeriesNew, returned.IsSeriesNew, "Series ID %s", expected.SeriesID)
	assert.Equal(t, expected.TMDBId, returned.TMDBId, "Series ID %s", expected.SeriesID)
	assert.Equal(t, expected.ProductionYear, returned.ProductionYear, "Series ID %s", expected.SeriesID)
}

func testReturnedSeasonIsCorrect(
	t *testing.T,
	expectedSeason SeasonItem,
	expectedSeasonID string,
	returnedSeason SeasonItem,
) {
	assert.Equal(t, expectedSeason.SeasonNumber, returnedSeason.SeasonNumber, "SeasonID %s", expectedSeasonID)
	assert.Equal(t, expectedSeason.Name, returnedSeason.Name, "SeasonID %s", expectedSeasonID)
	assert.InDelta(
		t,
		expectedSeason.AdditionDate.Unix(),
		returnedSeason.AdditionDate.Unix(),
		10,
		"SeasonID %s",
		expectedSeasonID,
	)
	assert.Equal(t, expectedSeason.IsSeasonNew, returnedSeason.IsSeasonNew, "SeasonID %s", expectedSeasonID)

	require.Len(
		t,
		returnedSeason.Episodes,
		len(expectedSeason.Episodes),
		"SeasonID %s. Data : %v",
		expectedSeasonID,
		returnedSeason,
	)
}

func testReturnedEpisodeIsCorrect(
	t *testing.T,
	expectedEpisodeItem EpisodeItem,
	expectedEpisodeID string,
	returnedEpisode EpisodeItem,
) {
	assert.Equal(t, expectedEpisodeItem.Name, returnedEpisode.Name, "episodeID %s", expectedEpisodeID)
	assert.InDelta(
		t,
		expectedEpisodeItem.AdditionDate.Unix(),
		returnedEpisode.AdditionDate.Unix(),
		10,
		"episodeID %s",
		expectedEpisodeID,
	)
	assert.Equal(t, expectedEpisodeItem.EpisodeNumber, returnedEpisode.EpisodeNumber, "episodeID %s", expectedEpisodeID)
}

func getSeriesItemBySeriesID(seriesID string, items *[]NewlyAddedSeriesItem) (*NewlyAddedSeriesItem, error) {
	for _, item := range *items {
		if item.SeriesID == seriesID {
			newItem := item
			return &newItem, nil
		}
	}
	return nil, errors.New("not found")
}

func assertLogsAreCorrect(t *testing.T, tt TableTests, recordedLogs *observer.ObservedLogs) {
	logs := recordedLogs.All()
	require.Len(t, logs, len(tt.loggedMessages))
	for i, log := range logs {
		assert.Equal(t, tt.loggedMessages[i].Message, log.Message)
		assert.ElementsMatch(t, tt.loggedMessages[i].Context, log.Context)
	}
}

func runGetNewlyAddedSeriesTest(t *testing.T, tt TableTests) {
	mockedApp, recordedLogs := testSeriesInitApp()
	mockedJellyfinBaseItem := tt.getSeriesBaseItems()
	expectedResult := tt.getExpectedResultFromBaseItem()

	mockItemsAPI := MockJellyfinItemsAPI{
		ExecuteGetAllItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			return &mockedJellyfinBaseItem, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}

	client := APIClient{
		ItemsAPI: mockItemsAPI,
	}
	returnedNewSeriesItems := client.GetNewlyAddedSeries(mockedApp)

	assertLogsAreCorrect(t, tt, recordedLogs)

	require.Len(t, *returnedNewSeriesItems, len(expectedResult))
	for _, expectedItem := range expectedResult {
		returnedItem, err := getSeriesItemBySeriesID(expectedItem.SeriesID, returnedNewSeriesItems)

		require.NoError(t, err, "Series ID %s", expectedItem.SeriesID)

		testReturnedSeriesIsCorrect(t, &expectedItem, returnedItem)

		if !expectedItem.IsSeriesNew {
			require.Len(t, returnedItem.NewSeasons, len(expectedItem.NewSeasons))
			for seasonID, season := range expectedItem.NewSeasons {
				testReturnedSeasonIsCorrect(t, season, seasonID, returnedItem.NewSeasons[seasonID])
				if !season.IsSeasonNew {
					for episodeID, episode := range season.Episodes {
						testReturnedEpisodeIsCorrect(
							t,
							episode,
							episodeID,
							returnedItem.NewSeasons[seasonID].Episodes[episodeID],
						)
					}
				}
			}
		}
	}
}

func getBaseItemIndexByID(id string) int {
	baseItems := getSeriesBaseItems()
	for i, item := range baseItems {
		if *item.Id == id {
			return i
		}
	}
	return 0
}

func getExpectedSeriesItemIndexByID(id string) int {
	expected := getExpectedResultFromBaseItem()
	for i, series := range expected {
		if series.SeriesID == id {
			return i
		}
	}
	return 0
}

func TestGetNewlyAddedSeries(t *testing.T) {
	tests := []TableTests{
		{
			name:                          "Valid data",
			getSeriesBaseItems:            getSeriesBaseItems,
			getExpectedResultFromBaseItem: getExpectedResultFromBaseItem,
		},
		{
			name:           "seriesName is null",
			loggedMessages: []observer.LoggedEntry{},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].Name = *jellyfinAPI.NewNullableString(nil)
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected[getExpectedSeriesItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].SeriesName = "Unknown"
				return expected
			},
		},
		{
			name:           "series productionYear is null",
			loggedMessages: []observer.LoggedEntry{},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].ProductionYear = *jellyfinAPI.NewNullableInt32(nil)
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected[getExpectedSeriesItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].ProductionYear = 0
				return expected
			},
		},
		{
			name: "series dateCreated is null",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "Found a series with no addition date. This can lead to inaccuracy when detecting newly added media.",
					},
					Context: []zapcore.Field{
						zap.String("Series ID", "1813f4b17e9d4a799641c09319b5ffcc"),
						zap.String("Series Name", "Series 1"),
					},
				},
			},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].DateCreated = *jellyfinAPI.NewNullableTime(nil)
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected[getExpectedSeriesItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].AdditionDate = time.Date(
					1970,
					01,
					01,
					00,
					00,
					00,
					00,
					time.UTC,
				)
				expected[getExpectedSeriesItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].IsSeriesNew = false
				// since the series has no creation date, the underlying season is taken
				expected[getExpectedSeriesItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].NewSeasons = map[string]SeasonItem{
					"f4971e32089041f3a3d6774277c2ccb9": {
						SeasonNumber: 1,
						Name:         "Season 1",
						AdditionDate: time.Now().AddDate(0, 0, -7),
						Episodes:     nil,
						IsSeasonNew:  true,
					},
				}
				return expected
			},
		},
		{
			name: "season dateCreated is null",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "Found a season with no addition date. This can lead to inaccuracy when detecting newly added media.",
					},
					Context: []zapcore.Field{
						zap.String("Series ID", "bb2222"),
						zap.String("Series Name", "Old Series 1"),
						zap.String("Season ID", "bb2222-s2"),
						zap.String("Season Name", "Season 2"),
					},
				},
			},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("bb2222-s2")].DateCreated = *jellyfinAPI.NewNullableTime(nil)
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				newSeason := expected[getExpectedSeriesItemIndexByID("bb2222")].NewSeasons["bb2222-s2"]
				newSeason.AdditionDate = time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)
				newSeason.IsSeasonNew = false
				// since the season has no creation date, the underlying episode is taken
				newSeason.Episodes = map[string]EpisodeItem{
					"bb2222-s2-e1": {
						EpisodeNumber: 1,
						Name:          "Episode 1",
						AdditionDate:  time.Now().AddDate(0, 0, -8),
					},
				}
				expected[getExpectedSeriesItemIndexByID("bb2222")].NewSeasons["bb2222-s2"] = newSeason
				return expected
			},
		},
		{
			name: "episode dateCreated is null",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "Found an episode with no addition date. This can lead to inaccuracy when detecting newly added media.",
					},
					Context: []zapcore.Field{
						zap.String("Episode ID", "cc3333-s1-e5"),
						zap.String("Episode Name", "Episode 5"),
						zap.String("Season Name", "Season 1"),
						zap.String("Season ID", "cc3333-s1"),
						zap.String("Series Name", "Very Old Series"),
						zap.String("Series ID", "cc3333"),
					},
				},
			},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("cc3333-s1-e5")].DateCreated = *jellyfinAPI.NewNullableTime(nil)
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected[getExpectedSeriesItemIndexByID("cc3333")].NewSeasons = map[string]SeasonItem{
					"cc3333-s2": expected[getExpectedSeriesItemIndexByID("cc3333")].NewSeasons["cc3333-s2"],
				}
				return expected
			},
		},
		{
			name:           "episode is virtual",
			loggedMessages: []observer.LoggedEntry{},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("cc3333-s1-e5")].LocationType = *jellyfinAPI.NewNullableLocationType(Ptr(jellyfinAPI.LOCATIONTYPE_VIRTUAL))
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected[getExpectedSeriesItemIndexByID("cc3333")].NewSeasons = map[string]SeasonItem{
					"cc3333-s2": expected[getExpectedSeriesItemIndexByID("cc3333")].NewSeasons["cc3333-s2"],
				}
				return expected
			},
		},
		{
			name:           "episode file location is nil",
			loggedMessages: []observer.LoggedEntry{},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("cc3333-s1-e5")].LocationType = *jellyfinAPI.NewNullableLocationType(nil)
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected[getExpectedSeriesItemIndexByID("cc3333")].NewSeasons = map[string]SeasonItem{
					"cc3333-s2": expected[getExpectedSeriesItemIndexByID("cc3333")].NewSeasons["cc3333-s2"],
				}
				return expected
			},
		},
		{
			name:           "series with no TMDBID",
			loggedMessages: []observer.LoggedEntry{},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems[getBaseItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].ProviderIds = map[string]string{
					"imdb": "1726",
				}
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected[getExpectedSeriesItemIndexByID("1813f4b17e9d4a799641c09319b5ffcc")].TMDBId = 0
				return expected
			},
		},
		{
			name: "Orphelin season",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "A season item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					},
					Context: []zapcore.Field{
						zap.String("Season ID", "bb2222-s1"),
						zap.String("Season Name", "Season 1"),
						zap.String("Not found Series Name", "Old Series 1"),
						zap.String("Not found Series ID", "bb2222"),
					},
				},
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "A season item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					},
					Context: []zapcore.Field{
						zap.String("Season ID", "bb2222-s2"),
						zap.String("Season Name", "Season 2"),
						zap.String("Not found Series Name", "Old Series 1"),
						zap.String("Not found Series ID", "bb2222"),
					},
				},
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "An episode item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					},
					Context: []zapcore.Field{
						zap.String("Expected Season ID", "bb2222-s2"),
						zap.String("Expected Season Name", "Season 2"),
						zap.String("Expected Series Name", "Old Series 1"),
						zap.String("Expected Series ID", "bb2222"),
						zap.String("Episode ID", "bb2222-s2-e1"),
						zap.String("Episode Name", "Episode 1"),
					},
				},
			},
			getSeriesBaseItems: func() []jellyfinAPI.BaseItemDto {
				baseItems := getSeriesBaseItems()
				baseItems = slices.Delete(baseItems, getBaseItemIndexByID("bb2222"), getBaseItemIndexByID("bb2222")+1)
				return baseItems
			},
			getExpectedResultFromBaseItem: func() []NewlyAddedSeriesItem {
				expected := getExpectedResultFromBaseItem()
				expected = slices.Delete(
					expected,
					getExpectedSeriesItemIndexByID("bb2222"),
					getExpectedSeriesItemIndexByID("bb2222")+1,
				)
				return expected
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runGetNewlyAddedSeriesTest(t, tt)
		})
	}
}
