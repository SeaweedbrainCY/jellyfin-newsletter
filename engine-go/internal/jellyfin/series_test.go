package jellyfin

import (
	"crypto/md5"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/clock"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type TableTests struct {
	getBaseItemsAndExpectedResults func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem)
	name                           string
	loggedMessages                 []observer.LoggedEntry
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
		Clock: clock.RealClock{},
	}, recordedLogs
}

func getBaseItemsAndExpectedResults() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
	observedDays := 30
	items := []struct {
		SeriesName                   string
		NumberOfEpisodes             int
		NumberOfEpisodesPerSeason    []int
		NumberOfNewEpisodesPerSeason []int
		IsSeriesNew                  bool
	}{
		{
			SeriesName:                   "Whole new series",
			IsSeriesNew:                  true,
			NumberOfEpisodesPerSeason:    []int{30, 30, 30},
			NumberOfNewEpisodesPerSeason: []int{30, 30, 30},
		},
		{
			SeriesName:                   "Whole new season",
			IsSeriesNew:                  false,
			NumberOfEpisodesPerSeason:    []int{30, 30, 30},
			NumberOfNewEpisodesPerSeason: []int{0, 0, 30},
		},
		{
			SeriesName:                   "Some new episodes",
			IsSeriesNew:                  false,
			NumberOfEpisodesPerSeason:    []int{30, 30, 30},
			NumberOfNewEpisodesPerSeason: []int{5, 0, 20},
		},
		{
			SeriesName:                   "No new episodes",
			IsSeriesNew:                  false,
			NumberOfEpisodesPerSeason:    []int{30, 30, 30},
			NumberOfNewEpisodesPerSeason: []int{0, 0, 0},
		},
		{
			// Every season has new episodes, but none of them is entirely new:
			// the series must not be considered as a whole new series.
			SeriesName:                   "Some new episodes in every season",
			IsSeriesNew:                  false,
			NumberOfEpisodesPerSeason:    []int{30, 30, 30},
			NumberOfNewEpisodesPerSeason: []int{5, 2, 20},
		},
		{
			// A single season series with only a few new episodes must not be
			// considered as a whole new series.
			SeriesName:                   "Some new episodes in a single season series",
			IsSeriesNew:                  false,
			NumberOfEpisodesPerSeason:    []int{30},
			NumberOfNewEpisodesPerSeason: []int{3},
		},
		{
			SeriesName:                   "Whole new single season series",
			IsSeriesNew:                  true,
			NumberOfEpisodesPerSeason:    []int{30},
			NumberOfNewEpisodesPerSeason: []int{30},
		},
	}

	baseItemDto := []jellyfinAPI.BaseItemDto{}
	expectedNewSeries := []NewlyAddedSeriesItem{}

	for _, item := range items {
		// Create the series
		seriesID := fmt.Sprintf("%x", md5.Sum([]byte(item.SeriesName)))
		baseItemDto = append(baseItemDto, jellyfinAPI.BaseItemDto{
			Id:             new(seriesID),
			Name:           *jellyfinAPI.NewNullableString(new(item.SeriesName)),
			ProductionYear: *jellyfinAPI.NewNullableInt32(new(int32(2023))),
			DateCreated:    *jellyfinAPI.NewNullableTime(new(time.Now().AddDate(0, 0, -7))),
			ProviderIds:    map[string]*string{"Tmdb": new("1027"), "Imdb": new("2276")},
			Type:           new(jellyfinAPI.BASEITEMKIND_SERIES),
			LocationType:   *jellyfinAPI.NewNullableLocationType(new(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
		})

		newSeries := NewlyAddedSeriesItem{
			SeriesName:     item.SeriesName,
			SeriesID:       seriesID,
			IsSeriesNew:    item.IsSeriesNew,
			TMDBId:         "1027",
			ProductionYear: 2023,
			AdditionDate:   time.Now().AddDate(0, 0, -7),
		}

		// Create the seasons
		newSeasons := map[string]SeasonItem{}
		for seasonNumber, numberOfEpisodes := range item.NumberOfEpisodesPerSeason {
			seasonName := item.SeriesName + " Season " + strconv.Itoa(seasonNumber)
			seasonID := fmt.Sprintf("%x", md5.Sum([]byte(seasonName)))
			baseItemDto = append(baseItemDto, jellyfinAPI.BaseItemDto{
				Id:             new(seasonID),
				Name:           *jellyfinAPI.NewNullableString(new(seasonName)),
				ProductionYear: *jellyfinAPI.NewNullableInt32(new(int32(2023))),
				DateCreated:    *jellyfinAPI.NewNullableTime(new(time.Now().AddDate(0, 0, -4))),
				ProviderIds:    map[string]*string{"Tmdb": new("1027"), "Imdb": new("2276")},
				Type:           new(jellyfinAPI.BASEITEMKIND_SEASON),
				LocationType:   *jellyfinAPI.NewNullableLocationType(new(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
				SeriesId:       *jellyfinAPI.NewNullableString(new(seriesID)),
				IndexNumber:    *jellyfinAPI.NewNullableInt32(new(int32(seasonNumber))),
			})

			// Create the episodes
			newEpisodes := map[string]EpisodeItem{}
			for episodeNumber := range numberOfEpisodes {
				episodeName := seasonName + " Episode " + strconv.Itoa(episodeNumber)
				episodeID := fmt.Sprintf("%x", md5.Sum([]byte(episodeName)))
				additionDate := time.Now().AddDate(0, 0, -2)
				if episodeNumber >= item.NumberOfNewEpisodesPerSeason[seasonNumber] {
					additionDate = time.Now().AddDate(0, 0, -50)
				}
				baseItemDto = append(baseItemDto, jellyfinAPI.BaseItemDto{
					Id:             new(episodeID),
					Name:           *jellyfinAPI.NewNullableString(new(episodeName)),
					DateCreated:    *jellyfinAPI.NewNullableTime(new(additionDate)),
					SeasonName:     *jellyfinAPI.NewNullableString(new(seasonName)),
					SeriesName:     *jellyfinAPI.NewNullableString(new(item.SeriesName)),
					Type:           new(jellyfinAPI.BASEITEMKIND_EPISODE),
					LocationType:   *jellyfinAPI.NewNullableLocationType(new(jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM)),
					SeriesId:       *jellyfinAPI.NewNullableString(new(seriesID)),
					SeasonId:       *jellyfinAPI.NewNullableString(new(seasonID)),
					ProductionYear: *jellyfinAPI.NewNullableInt32(new(int32(2023))),
					IndexNumber:    *jellyfinAPI.NewNullableInt32(new(int32(episodeNumber))),
				})
				if additionDate.After(time.Now().AddDate(0, 0, observedDays*-1)) {
					newEpisodes[episodeID] = EpisodeItem{
						Name:          episodeName,
						AdditionDate:  additionDate,
						EpisodeNumber: int32(episodeNumber),
					}
				}
			}
			if len(newEpisodes) != 0 {
				newSeasons[seasonID] = SeasonItem{
					SeasonNumber: int32(seasonNumber),
					Name:         seasonName,
					AdditionDate: time.Now().AddDate(0, 0, -4),
					IsSeasonNew:  numberOfEpisodes == item.NumberOfNewEpisodesPerSeason[seasonNumber],
					Episodes:     newEpisodes,
				}
			}
		}
		if len(newSeasons) != 0 {
			newSeries.NewSeasons = newSeasons
			expectedNewSeries = append(expectedNewSeries, newSeries)
		}
	}
	return baseItemDto, expectedNewSeries
}

// cloneTestData returns copies of the shared fixtures so a test case can mutate
// them without leaking changes into other test cases.
// Note: BaseItemDto elements are shallow-copied. Replacing a field is safe, but
// mutating through a pointer/map field (e.g. *item.Id = "x" or item.ProviderIds["k"] = v)
// would still affect the original.
func cloneTestData(
	baseItems []jellyfinAPI.BaseItemDto,
	expectedResults []NewlyAddedSeriesItem,
) ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
	baseItemsCopy := slices.Clone(baseItems)
	expectedCopy := slices.Clone(expectedResults)
	for i := range expectedCopy {
		seasons := maps.Clone(expectedCopy[i].NewSeasons)
		for id, season := range seasons {
			season.Episodes = maps.Clone(season.Episodes)
			seasons[id] = season
		}
		expectedCopy[i].NewSeasons = seasons
	}
	return baseItemsCopy, expectedCopy
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
	mockedJellyfinBaseItem, expectedResult := tt.getBaseItemsAndExpectedResults()

	mockLibraryAPI := MockJellyfinLibraryAPI{
		ExecuteGetAllItemsByFolderID: func() (*[]jellyfinAPI.BaseItemDto, error) {
			return &mockedJellyfinBaseItem, nil
		},
		ExecuteGetRootFolderIDByName: func() (string, error) {
			return "id", nil
		},
	}

	client := APIClient{
		LibraryAPI: mockLibraryAPI,
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

func getBaseItemIndexByName(baseItems []jellyfinAPI.BaseItemDto, name string) int {
	for i, item := range baseItems {
		if item.Name.Get() != nil && *item.Name.Get() == name {
			return i
		}
	}
	return 0
}

func getExpectedSeriesIDIndexByName(expectedResults []NewlyAddedSeriesItem, name string) int {
	for i, series := range expectedResults {
		if series.SeriesName == name {
			return i
		}
	}
	return 0
}

func getExpectedSeriesItemBySeriesName(expectedResults []NewlyAddedSeriesItem, name string) NewlyAddedSeriesItem {
	for _, series := range expectedResults {
		if series.SeriesName == name {
			return series
		}
	}
	return NewlyAddedSeriesItem{}
}

// idFromName rebuilds the ID generated by the fixtures for a given item name.
func idFromName(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}

// removeExpectedEpisode drops an episode from the expected results, for test cases
// where the episode is supposed to be ignored by the parser.
func removeExpectedEpisode(
	expectedResults []NewlyAddedSeriesItem,
	seriesName string,
	seasonName string,
	episodeName string,
) {
	seriesIndex := getExpectedSeriesIDIndexByName(expectedResults, seriesName)
	delete(expectedResults[seriesIndex].NewSeasons[idFromName(seasonName)].Episodes, idFromName(episodeName))
}

func TestGetNewlyAddedSeries(t *testing.T) {
	refBaseItems, refExpectedResults := getBaseItemsAndExpectedResults()
	tests := []TableTests{
		{
			name:                           "Valid data",
			getBaseItemsAndExpectedResults: getBaseItemsAndExpectedResults,
		},
		{
			name:           "seriesName is null",
			loggedMessages: []observer.LoggedEntry{},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				baseItems[getBaseItemIndexByName(baseItems, "Whole new series")].Name = *jellyfinAPI.NewNullableString(nil)
				expectedResults[getExpectedSeriesIDIndexByName(expectedResults, "Whole new series")].SeriesName = ""
				return baseItems, expectedResults
			},
		},
		{
			name:           "series productionYear is null",
			loggedMessages: []observer.LoggedEntry{},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				baseItems[getBaseItemIndexByName(baseItems, "Whole new series")].ProductionYear = *jellyfinAPI.NewNullableInt32(nil)
				expectedResults[getExpectedSeriesIDIndexByName(expectedResults, "Whole new series")].ProductionYear = 0
				return baseItems, expectedResults
			},
		},
		{
			name: "series dateCreated is null",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "Found a series with no addition date. This can lead to inaccuracies.",
					},
					Context: []zapcore.Field{
						zap.String(
							"Series ID",
							getExpectedSeriesItemBySeriesName(refExpectedResults, "Whole new series").SeriesID,
						),
						zap.String("Series Name", "Whole new series"),
					},
				},
			},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				baseItems[getBaseItemIndexByName(baseItems, "Whole new series")].DateCreated = *jellyfinAPI.NewNullableTime(nil)
				expectedResults[getExpectedSeriesIDIndexByName(expectedResults, "Whole new series")].AdditionDate = time.Date(
					1970,
					01,
					01,
					00,
					00,
					00,
					00,
					time.UTC,
				)
				return baseItems, expectedResults
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
						zap.String(
							"Episode ID",
							fmt.Sprintf("%x", md5.Sum([]byte("Whole new series Season 1 Episode 1"))),
						),
						zap.String("Episode Name", "Whole new series Season 1 Episode 1"),
						zap.String("Season Name", "Whole new series Season 1"),
						zap.String("Season ID", fmt.Sprintf("%x", md5.Sum([]byte("Whole new series Season 1")))),
						zap.String("Series Name", "Whole new series"),
						zap.String("Series ID", fmt.Sprintf("%x", md5.Sum([]byte("Whole new series")))),
					},
				},
			},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				baseItems[getBaseItemIndexByName(baseItems, "Whole new series Season 1 Episode 1")].DateCreated = *jellyfinAPI.NewNullableTime(nil)
				seriesIndex := getExpectedSeriesIDIndexByName(expectedResults, "Whole new series")

				// Without an addition date the episode is not detected as new, so its
				// season is not entirely new anymore, and neither is the series.
				removeExpectedEpisode(
					expectedResults,
					"Whole new series",
					"Whole new series Season 1",
					"Whole new series Season 1 Episode 1",
				)
				seasonID := idFromName("Whole new series Season 1")
				season := expectedResults[seriesIndex].NewSeasons[seasonID]
				season.IsSeasonNew = false
				expectedResults[seriesIndex].NewSeasons[seasonID] = season
				expectedResults[seriesIndex].IsSeriesNew = false

				return baseItems, expectedResults
			},
		},
		{
			name:           "episode is virtual",
			loggedMessages: []observer.LoggedEntry{},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				episodeIndex := getBaseItemIndexByName(baseItems, "Some new episodes Season 0 Episode 1")
				baseItems[episodeIndex].LocationType = *jellyfinAPI.NewNullableLocationType(
					new(jellyfinAPI.LOCATIONTYPE_VIRTUAL),
				)
				removeExpectedEpisode(
					expectedResults,
					"Some new episodes",
					"Some new episodes Season 0",
					"Some new episodes Season 0 Episode 1",
				)
				return baseItems, expectedResults
			},
		},
		{
			name:           "episode locationType is null",
			loggedMessages: []observer.LoggedEntry{},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				episodeIndex := getBaseItemIndexByName(baseItems, "Some new episodes Season 0 Episode 1")
				baseItems[episodeIndex].LocationType = *jellyfinAPI.NewNullableLocationType(nil)
				removeExpectedEpisode(
					expectedResults,
					"Some new episodes",
					"Some new episodes Season 0",
					"Some new episodes Season 0 Episode 1",
				)
				return baseItems, expectedResults
			},
		},
		{
			name:           "series has no TMDB ID",
			loggedMessages: []observer.LoggedEntry{},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				baseItems[getBaseItemIndexByName(baseItems, "Whole new series")].ProviderIds = map[string]*string{
					"Imdb": new("2276"),
				}
				expectedResults[getExpectedSeriesIDIndexByName(expectedResults, "Whole new series")].TMDBId = ""
				return baseItems, expectedResults
			},
		},
		{
			name:           "series TMDB ID is null",
			loggedMessages: []observer.LoggedEntry{},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				baseItems[getBaseItemIndexByName(baseItems, "Whole new series")].ProviderIds = map[string]*string{
					"Tmdb": nil,
					"Imdb": new("2276"),
				}
				expectedResults[getExpectedSeriesIDIndexByName(expectedResults, "Whole new series")].TMDBId = ""
				return baseItems, expectedResults
			},
		},
		{
			name: "season belongs to an unknown series",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "A season item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					},
					Context: []zapcore.Field{
						zap.String("Season ID", idFromName("Orphan season")),
						zap.String("Season Name", "Orphan season"),
						zap.String("Not found Series Name", "Unknown series"),
						zap.String("Not found Series ID", idFromName("Unknown series")),
					},
				},
			},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				seasonID := idFromName("Orphan season")
				seasonName := "Orphan season"
				unknownSeriesID := idFromName("Unknown series")
				unknownSeriesName := "Unknown series"
				baseItems = append(baseItems, jellyfinAPI.BaseItemDto{
					Id:           new(seasonID),
					Name:         *jellyfinAPI.NewNullableString(new(seasonName)),
					DateCreated:  *jellyfinAPI.NewNullableTime(new(time.Now().AddDate(0, 0, -4))),
					Type:         new(jellyfinAPI.BASEITEMKIND_SEASON),
					LocationType: *jellyfinAPI.NewNullableLocationType(new(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
					SeriesId:     *jellyfinAPI.NewNullableString(new(unknownSeriesID)),
					SeriesName:   *jellyfinAPI.NewNullableString(new(unknownSeriesName)),
					IndexNumber:  *jellyfinAPI.NewNullableInt32(new(int32(1))),
				})
				return baseItems, expectedResults
			},
		},
		{
			name: "season has no series ID",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "A season item is ignored because it has no series ID.",
					},
					Context: []zapcore.Field{
						zap.String("Season ID", idFromName("Season without series")),
					},
				},
			},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				seasonID := idFromName("Season without series")
				seasonName := "Season without series"
				baseItems = append(baseItems, jellyfinAPI.BaseItemDto{
					Id:           new(seasonID),
					Name:         *jellyfinAPI.NewNullableString(new(seasonName)),
					DateCreated:  *jellyfinAPI.NewNullableTime(new(time.Now().AddDate(0, 0, -4))),
					Type:         new(jellyfinAPI.BASEITEMKIND_SEASON),
					LocationType: *jellyfinAPI.NewNullableLocationType(new(jellyfinAPI.LOCATIONTYPE_VIRTUAL)),
					SeriesId:     *jellyfinAPI.NewNullableString(nil),
					IndexNumber:  *jellyfinAPI.NewNullableInt32(new(int32(1))),
				})
				return baseItems, expectedResults
			},
		},
		{
			name: "episode has no season ID",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "An episode item is ignored because it has no series ID or season ID.",
					},
					Context: []zapcore.Field{
						zap.String("Episode ID", idFromName("Some new episodes Season 0 Episode 1")),
						zap.String("Episode Name", "Some new episodes Season 0 Episode 1"),
						zap.String("Expected Series Name", "Some new episodes"),
						zap.String("Expected Series ID", idFromName("Some new episodes")),
						zap.String("Expected Season Name", "Some new episodes Season 0"),
						zap.String("Expected Season ID", ""),
					},
				},
			},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				episodeIndex := getBaseItemIndexByName(baseItems, "Some new episodes Season 0 Episode 1")
				baseItems[episodeIndex].SeasonId = *jellyfinAPI.NewNullableString(nil)
				removeExpectedEpisode(
					expectedResults,
					"Some new episodes",
					"Some new episodes Season 0",
					"Some new episodes Season 0 Episode 1",
				)
				return baseItems, expectedResults
			},
		},
		{
			name: "episode has no series ID",
			loggedMessages: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zapcore.WarnLevel,
						Message: "An episode item is ignored because it has no series ID or season ID.",
					},
					Context: []zapcore.Field{
						zap.String("Episode ID", idFromName("Some new episodes Season 0 Episode 1")),
						zap.String("Episode Name", "Some new episodes Season 0 Episode 1"),
						zap.String("Expected Series Name", "Some new episodes"),
						zap.String("Expected Series ID", ""),
						zap.String("Expected Season Name", "Some new episodes Season 0"),
						zap.String("Expected Season ID", idFromName("Some new episodes Season 0")),
					},
				},
			},
			getBaseItemsAndExpectedResults: func() ([]jellyfinAPI.BaseItemDto, []NewlyAddedSeriesItem) {
				baseItems, expectedResults := cloneTestData(refBaseItems, refExpectedResults)
				episodeIndex := getBaseItemIndexByName(baseItems, "Some new episodes Season 0 Episode 1")
				baseItems[episodeIndex].SeriesId = *jellyfinAPI.NewNullableString(nil)
				removeExpectedEpisode(
					expectedResults,
					"Some new episodes",
					"Some new episodes Season 0",
					"Some new episodes Season 0 Episode 1",
				)
				return baseItems, expectedResults
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runGetNewlyAddedSeriesTest(t, tt)
		})
	}
}
