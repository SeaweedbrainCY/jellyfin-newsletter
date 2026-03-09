package template

import (
	"strconv"
	"testing"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/i18n"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func timePtr(t time.Time) *time.Time {
	return &t
}

func getJellyfinNewMovies() []jellyfin.MovieItem {
	return []jellyfin.MovieItem{
		{
			ID:             "7dcf7149-f710-46d5-a50c-626e3486259b",
			Name:           "Star Wars: Episode II - Attack of the Clones",
			AdditionDate:   timePtr(time.Date(2026, 01, 01, 01, 01, 0, 0, time.UTC)),
			TMDBId:         "1274",
			ProductionYear: int32(2025),
			Overview:       "Following an assassination attempt on Senator Padmé Amidala, Jedi Knights Anakin Skywalker and Obi-Wan Kenobi investigate a mysterious plot into the heart of the Separatist movement and the beginning of the Clone Wars.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/oZNPzxqM2s5DyVWab09NTQScDQt.jpg",
		},
		{
			ID:             "fd9416da-9026-4219-95b4-0dae418d2b5d",
			Name:           "Oppenheimer",
			AdditionDate:   timePtr(time.Date(2026, 01, 02, 01, 01, 0, 0, time.UTC)),
			TMDBId:         "1273",
			ProductionYear: int32(2023),
			Overview:       "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
		},
	}
}

func getJellyfinNewSeriesItems() []jellyfin.NewlyAddedSeriesItem {
	return []jellyfin.NewlyAddedSeriesItem{
		{
			SeriesName:     "Game of thrones",
			SeriesID:       "3d7b0576-370c-48d3-b7c3-7c49f612afc9",
			IsSeriesNew:    true,
			NewSeasons:     nil,
			TMDBId:         "1735",
			ProductionYear: 2021,
			AdditionDate:   time.Date(2026, 01, 03, 01, 01, 0, 0, time.UTC),
			Overview:       "Seven noble families fight for control of the mythical land of Westeros. Friction between the houses leads to full-scale war. All while a very ancient evil awakens in the farthest north. Amidst the war, a neglected military order of misfits, the Night's Watch, is all that stands between the realms of men and icy horrors beyond.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/1XS1oqL89opfnbLl8WnZY1O1uJx.jpg",
		},
		{
			SeriesName:  "How I Met Your Mother",
			SeriesID:    "c828b892-64f8-4def-88b7-dc3d9072a147",
			IsSeriesNew: false,
			NewSeasons: map[string]jellyfin.SeasonItem{
				"eb122458-7697-4151-80bf-8cb911685398": {
					SeasonNumber: int32(9),
					Name:         "Season 9",
					AdditionDate: time.Date(2026, 01, 06, 01, 01, 0, 0, time.UTC),
					Episodes:     nil,
					IsSeasonNew:  true,
				},
			},
			TMDBId:         "9372",
			ProductionYear: 2016,
			AdditionDate:   time.Date(2025, 10, 01, 01, 01, 0, 0, time.UTC),
			Overview:       "A father recounts to his children - through a series of flashbacks - the journey he and his four best friends took leading up to him meeting their mother.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/b34jPzmB0wZy7EjUZoleXOl2RRI.jpg",
		},
		{
			SeriesName:  "Stranger Things",
			SeriesID:    "66b82dc8-d655-44b3-b8d7-5a5eaaa581cd",
			IsSeriesNew: false,
			NewSeasons: map[string]jellyfin.SeasonItem{
				"d4da2014-b39a-4d2f-9f12-b45492af4ae6": {
					SeasonNumber: int32(1),
					Name:         "Season 1",
					AdditionDate: time.Date(2025, 12, 05, 01, 01, 0, 0, time.UTC),
					Episodes: map[string]jellyfin.EpisodeItem{
						"a9987481-6aa8-4cd7-850b-dd962c235685": {
							Name:          "Episode 5",
							AdditionDate:  time.Date(2026, 01, 04, 01, 01, 0, 0, time.UTC),
							EpisodeNumber: int32(5),
						},
					},
					IsSeasonNew: false,
				},
				"a2b5bd7f-7c98-4f35-b27b-2c6901982ec1": {
					SeasonNumber: int32(2),
					Name:         "Season 2",
					AdditionDate: time.Date(2025, 12, 05, 01, 01, 0, 0, time.UTC),
					Episodes: map[string]jellyfin.EpisodeItem{
						"a915c746-c2dc-4f5e-a9e7-f6a4a8595d2e": {
							Name:          "Episode 10",
							AdditionDate:  time.Date(2026, 01, 02, 01, 01, 0, 0, time.UTC),
							EpisodeNumber: int32(10),
						},
					},
					IsSeasonNew: false,
				},
			},
			TMDBId:         "3628",
			ProductionYear: 2013,
			AdditionDate:   time.Date(2025, 11, 01, 01, 01, 0, 0, time.UTC),
			Overview:       "When a young boy vanishes, a small town uncovers a mystery involving secret experiments, terrifying supernatural forces, and one strange little girl.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/uOOtwVbSr4QDjAGIifLDwpb2Pdl.jpg",
		},
		{
			SeriesName:  "Family Guy",
			SeriesID:    "94327c53-7a32-4e7c-84f4-5b1a6e71dd35",
			IsSeriesNew: false,
			NewSeasons: map[string]jellyfin.SeasonItem{
				"4d747e44-b807-4360-a308-62c098026e6f": {
					SeasonNumber: int32(24),
					Name:         "Season 24",
					AdditionDate: time.Date(2025, 12, 05, 01, 01, 0, 0, time.UTC),
					Episodes: map[string]jellyfin.EpisodeItem{
						"1c50012a-8c35-4c31-9225-f063f563b83e": {
							Name:          "Episode 1",
							AdditionDate:  time.Date(2026, 01, 05, 01, 01, 0, 0, time.UTC),
							EpisodeNumber: int32(1),
						},
						"3140d805-0c1b-4268-9ee1-1acca7ba565a": {
							Name:          "Episode 2",
							AdditionDate:  time.Date(2026, 01, 04, 01, 01, 0, 0, time.UTC),
							EpisodeNumber: int32(2),
						},
						"48558cb6-2b63-4a9f-89f0-e9369bd751e3": {
							Name:          "Episode 7",
							AdditionDate:  time.Date(2026, 01, 05, 01, 01, 0, 0, time.UTC),
							EpisodeNumber: int32(7),
						},
						"5dec0e7c-11a4-459e-ae05-a292810851bf": {
							Name:          "Episode 3",
							AdditionDate:  time.Date(2026, 01, 05, 01, 01, 0, 0, time.UTC),
							EpisodeNumber: int32(3),
						},
					},
					IsSeasonNew: false,
				},
			},
			TMDBId:         "3271",
			ProductionYear: 2015,
			AdditionDate:   time.Date(2025, 11, 01, 01, 01, 0, 0, time.UTC),
			Overview:       "Sick, twisted, politically incorrect and Freakin' Sweet animated series featuring the adventures of the dysfunctional Griffin family. Bumbling Peter and long-suffering Lois have three kids. Stewie (a brilliant but sadistic baby bent on killing his mother and taking over the world), Meg (the oldest, and is the most unpopular girl in town) and Chris (the middle kid, he's not very bright but has a passion for movies). The final member of the family is Brian - a talking dog and much more than a pet, he keeps Stewie in check whilst sipping Martinis and sorting through his own life issues.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/3PFsEuAiyLkWsP4GG6dIV37Q6gu.jpg",
		},
	}
}

func getExpectedNewMediaTemplateData() newMediaTemplateData {
	title := "New items from " + time.Now().
		AddDate(0, 0, -30).
		Format("Monday January 2006 2006-01-02") +
		" to " + time.Now().
		Format("Monday January 2006 2006-01-02")
	subtitle := "Subtitle: " + title

	newMovies := []newMovieItemTemplateData{
		{
			PosterURL:    "https://image.tmdb.org/t/p/w500/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
			Name:         "Oppenheimer",
			AdditionDate: "2026-01-02",
			Overview:     "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
			AddedOnLabel: "Added on",
		},
		{
			PosterURL:    "https://image.tmdb.org/t/p/w500/oZNPzxqM2s5DyVWab09NTQScDQt.jpg",
			Name:         "Star Wars: Episode II - Attack of the Clones",
			AdditionDate: "2026-01-01",
			Overview:     "Following an assassination attempt on Senator Padmé Amidala, Jedi Knights Anakin Skywalker and Obi-Wan Kenobi investigate a mysterious plot into the heart of the Separatist movement and the beginning of the Clone Wars.",
			AddedOnLabel: "Added on",
		},
	}

	newSeries := []newSeriesItemTemplateData{
		{
			// Old series, new episodes in 1 season
			PosterURL:      "https://image.tmdb.org/t/p/w500/3PFsEuAiyLkWsP4GG6dIV37Q6gu.jpg",
			SeriesName:     "Family Guy",
			AdditionDate:   "2026-01-05",
			Overview:       "Sick, twisted, politically incorrect and Freakin' Sweet animated series featuring the adventures of the dysfunctional Griffin family. Bumbling Peter and long-suffering Lois have three kids. Stewie (a brilliant but sadistic baby bent on killing his mother and taking over the world), Meg (the oldest, and is the most unpopular girl in town) and Chris (the middle kid, he's not very bright but has a passion for movies). The final member of the family is Brian - a talking dog and much more than a pet, he keeps Stewie in check whilst sipping Martinis and sorting through his own life issues.",
			AddedOnLabel:   "Added on",
			NewSeriesTitle: "Family Guy: Season 24, Episodes 1-3 & 7",
		},
		{
			// Whole new series
			PosterURL:      "https://image.tmdb.org/t/p/w500/1XS1oqL89opfnbLl8WnZY1O1uJx.jpg",
			SeriesName:     "Game of thrones",
			AdditionDate:   "2026-01-03",
			Overview:       "Seven noble families fight for control of the mythical land of Westeros. Friction between the houses leads to full-scale war. All while a very ancient evil awakens in the farthest north. Amidst the war, a neglected military order of misfits, the Night's Watch, is all that stands between the realms of men and icy horrors beyond.",
			AddedOnLabel:   "Added on",
			NewSeriesTitle: "Game of thrones: Season 1-8",
		},
		{
			// Old series, new season
			PosterURL:      "https://image.tmdb.org/t/p/w500/b34jPzmB0wZy7EjUZoleXOl2RRI.jpg",
			SeriesName:     "How I Met Your Mother",
			AdditionDate:   "2026-01-06",
			Overview:       "A father recounts to his children - through a series of flashbacks - the journey he and his four best friends took leading up to him meeting their mother.",
			AddedOnLabel:   "Added on",
			NewSeriesTitle: "How I Met Your Mother: Season 9",
		},
		{
			// Old series, new episodes in 2 seasons
			PosterURL:      "https://image.tmdb.org/t/p/w500/uOOtwVbSr4QDjAGIifLDwpb2Pdl.jpg",
			SeriesName:     "Stranger Things",
			AdditionDate:   "2026-01-04",
			Overview:       "When a young boy vanishes, a small town uncovers a mystery involving secret experiments, terrifying supernatural forces, and one strange little girl.",
			AddedOnLabel:   "Added on",
			NewSeriesTitle: "Stranger Things: Season 1-2",
		},
	}

	return newMediaTemplateData{
		HTMLLang:                     "en",
		HTMLDir:                      "ltr",
		Title:                        title,
		Subtitle:                     subtitle,
		JellyfinURL:                  "https://jellyfin.example.com",
		DiscoverNowLabel:             "Discover now",
		DisplayNewMovies:             true,
		NewFilmLabel:                 "New movies:",
		NewMovies:                    newMovies,
		DisplayNewSeries:             true,
		NewSeriesLabel:               "New shows:",
		NewSeries:                    newSeries,
		CurrentlyAvailableLabel:      "Currently available in Jellyfin:",
		MoviesCount:                  strconv.Itoa(54),
		SeriesCount:                  strconv.Itoa(1253),
		MoviesLabel:                  "Movies",
		SeriesLabel:                  "Series",
		FooterLabel:                  "You are recieving this email because you are using seaweedbrain's Jellyfin server. If you want to stop receiving these emails, you can unsubscribe by notifying stop@example.com.",
		FooterProjectLinkLabel:       "Jellyfin Newsletter",
		FooterOpenSourceProjectLabel: "is an open source project.",
		FooterDevelopedByLabel:       "Developed with ❤️ by <a href=\"https://github.com/SeaweedbrainCY/\" class=\"footer-link\">SeaweedbrainCY</a> and <a href=\"https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors\" class=\"footer-link\">contributors</a>.",
		FooterLicenceAndCopyright:    "Copyright © 2025 Nathan Stchepinsky, licensed under AGPLv3.",
	}
}

func getAppContext() (*app.ApplicationContext, *observer.ObservedLogs) {
	localizer, _ := i18n.NewLocalizer("en")
	loggerCore, recordedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(loggerCore)
	return &app.ApplicationContext{
		Localizer: localizer,
		Logger:    logger,
		Config: &config.Configuration{
			Jellyfin: config.JellyfinConfig{
				ObservedPeriodDays: 30,
			},
			EmailTemplate: config.EmailTemplateConfig{
				Theme:                   "classic",
				Language:                "en",
				Subject:                 "Not important",
				Title:                   "New items from {{.StartDayName}} {{.StartMonthName}} {{.StartYear}} {{.StartDate}} to {{.DayName}} {{.MonthName}} {{.Year}} {{.Date}}",
				Subtitle:                "Subtitle: New items from {{.StartDayName}} {{.StartMonthName}} {{.StartYear}} {{.StartDate}} to {{.DayName}} {{.MonthName}} {{.Year}} {{.Date}}",
				JellyfinURL:             "https://jellyfin.example.com",
				UnsubscribeEmail:        "stop@example.com",
				JellyfinOwnerName:       "seaweedbrain",
				DisplayOverviewMaxItems: 0,
			},
		},
	}, recordedLogs
}

func TestBuildNewMediaTemplateData(t *testing.T) {
	tests := []struct {
		name                                string
		getAppContextFunc                   func() (*app.ApplicationContext, *observer.ObservedLogs)
		getExpectedNewMediaTemplateDataFunc func() newMediaTemplateData
		getJellyfinNewMovies                func() []jellyfin.MovieItem
		getJellyfinNewSeriesItems           func() []jellyfin.NewlyAddedSeriesItem
		movieCount                          int
		episodeCount                        int
	}{
		{
			name:                                "Valid data",
			getAppContextFunc:                   getAppContext,
			getExpectedNewMediaTemplateDataFunc: getExpectedNewMediaTemplateData,
			getJellyfinNewMovies:                getJellyfinNewMovies,
			getJellyfinNewSeriesItems:           getJellyfinNewSeriesItems,
			movieCount:                          54,
			episodeCount:                        1253,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, recordedLogs := test.getAppContextFunc()
			expectedTemplateData := test.getExpectedNewMediaTemplateDataFunc()
			newSeries := test.getJellyfinNewSeriesItems()
			newMovies := test.getJellyfinNewMovies()
			templateData, err := buildNewMediaTemplateData(
				&newMovies,
				&newSeries,
				int32(test.movieCount),
				int32(test.episodeCount),
				app,
			)
			require.NoError(t,
				err)
			assert.Empty(t, recordedLogs)
			assert.Equal(t, expectedTemplateData, templateData)
		})
	}
}

//func TestCheckIfThemeIsAvailable(t *testing.T) {
//
//}
