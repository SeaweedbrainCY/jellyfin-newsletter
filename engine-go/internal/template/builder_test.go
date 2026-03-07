package template

import (
	"strconv"
	"testing"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
)

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
			AdditionDate: "2026-01-01",
			Overview:     "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
			AddedOnLabel: "Added on",
		},
		{
			PosterURL:    "https://image.tmdb.org/t/p/w500/oZNPzxqM2s5DyVWab09NTQScDQt.jpg",
			Name:         "Star Wars: Episode II - Attack of the Clones",
			AdditionDate: "2026-01-02",
			Overview:     "Following an assassination attempt on Senator Padmé Amidala, Jedi Knights Anakin Skywalker and Obi-Wan Kenobi investigate a mysterious plot into the heart of the Separatist movement and the beginning of the Clone Wars.",
			AddedOnLabel: "Added on",
		},
	}

	newSeries := []newSeriesItemTemplateData{
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
			AdditionDate:   "2026-01-04",
			Overview:       "A father recounts to his children - through a series of flashbacks - the journey he and his four best friends took leading up to him meeting their mother.",
			AddedOnLabel:   "Added on",
			NewSeriesTitle: "How I Met Your Mother: Season 9",
		},
		{
			// Old series, new episodes in 1 season
			PosterURL:      "https://image.tmdb.org/t/p/w500/3PFsEuAiyLkWsP4GG6dIV37Q6gu.jpg",
			SeriesName:     "Family Guy",
			AdditionDate:   "2026-01-05",
			Overview:       "Sick, twisted, politically incorrect and Freakin' Sweet animated series featuring the adventures of the dysfunctional Griffin family. Bumbling Peter and long-suffering Lois have three kids. Stewie (a brilliant but sadistic baby bent on killing his mother and taking over the world), Meg (the oldest, and is the most unpopular girl in town) and Chris (the middle kid, he's not very bright but has a passion for movies). The final member of the family is Brian - a talking dog and much more than a pet, he keeps Stewie in check whilst sipping Martinis and sorting through his own life issues.",
			AddedOnLabel:   "Added on",
			NewSeriesTitle: "Family Guy: Season 24, Episodes 1-5 & 7",
		},
		{
			// Old series, new episodes in 2 seasons
			PosterURL:      "https://image.tmdb.org/t/p/w500/uOOtwVbSr4QDjAGIifLDwpb2Pdl.jpg",
			SeriesName:     "Stranger Things",
			AdditionDate:   "2026-01-06",
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

func getAppContext() *app.ApplicationContext {
	return &app.ApplicationContext{
		Config: &config.Configuration{
			EmailTemplate: config.EmailTemplateConfig{
				Theme:                   "classic",
				Language:                "en",
				Subject:                 "Not important",
				Title:                   "New items from {{.StartDayName}} {{.StartMonthName}} {{.StartDate}}",
				Subtitle:                "Subtitle: New items from {{.StartDayName}} {{.StartMonthName}} {{.StartDate}}",
				JellyfinURL:             "https://jellyfin.example.com",
				UnsubscribeEmail:        "stop@example.com",
				JellyfinOwnerName:       "seaweedbrain",
				DisplayOverviewMaxItems: 0,
			},
		},
	}
}

func TestBuildNewMediaTemplateData(t *testing.T) {

}

func TestCheckIfThemeIsAvailable(t *testing.T) {

}
