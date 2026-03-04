package template

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"go.uber.org/zap"
)

type newMovieItemTemplateData struct {
	PosterURL    string
	Name         string
	AddedOnLabel string
	AdditionDate string
	Overview     string
}
type newSeriesItemTemplateData struct {
	PosterURL      string
	SeriesName     string
	AddedOnLabel   string
	AdditionDate   string
	Overview       string
	NewSeriesTitle string
}

type newMediaTemplateData struct {
	HTMLLang                     string
	HTMLDir                      string
	Title                        string
	Subtitle                     string
	JellyfinURL                  string
	DiscoverNowLabel             string
	DisplayNewMovies             bool
	NewFilmLabel                 string
	NewMovies                    []newMovieItemTemplateData
	DisplayNewSeries             bool
	NewSeriesLabel               string
	NewSeries                    []newSeriesItemTemplateData
	CurrentlyAvailableLabel      string
	MoviesCount                  string
	MoviesLabel                  string
	SeriesCount                  string
	SeriesLabel                  string
	FooterLabel                  string
	FooterProjectLinkLabel       string
	FooterOpenSourceProjectLabel string
	FooterDevelopedByLabel       string
	FooterLicenceAndCopyright    string
}

type titlePlaceholders struct {
	Date             string
	DayName          string
	DayNumber        string
	MonthName        string
	MonthNumber      string
	Year             string
	StartDate        string
	StartDayName     string
	StartDayNumber   string
	StartMonthName   string
	StartMonthNumber string
	StartYear        string
}

func getNewMediaHTMLTemplate(themes_type string, app *app.ApplicationContext) (*template.Template, error) {
	filePath := filepath.Join("themes", themes_type, app.Config.EmailTemplate.Theme)
	tmpl, err := template.ParseFiles(filePath)

	if err != nil {
		app.Logger.Error(
			"Impossible to open the template file. Email HTML building will fail.",
			zap.String("filePath", filePath),
			zap.Error(err),
		)
		return "", err
	}

	return string(template), nil
}

// compressNumbers converts [1 2 3 6 7 9] → "1-3 & 6-7 & 9".
func compressNumbers(nums []int) string {
	if len(nums) == 0 {
		return ""
	}

	sort.Ints(nums)

	var parts []string
	start := nums[0]
	prev := nums[0]

	for i := 1; i < len(nums); i++ {
		if nums[i] == prev+1 {
			prev = nums[i]
			continue
		}

		parts = append(parts, formatRange(start, prev))
		start = nums[i]
		prev = nums[i]
	}

	compressedString := strings.Join(parts, ", ")

	lastPart := formatRange(start, prev)

	if compressedString != "" {
		compressedString = compressedString + " & " + lastPart
	} else {
		compressedString = lastPart
	}

	return compressedString
}

func formatRange(start, end int) string {
	if start == end {
		return fmt.Sprintf("%d", start)
	}
	return fmt.Sprintf("%d-%d", start, end)
}

// New series item are not only displayed with the series name.
// Sometime the whole series is new, sometimes it's only seasons or episodes
// This func build the 'title' with will be displayed in the template and
// which depends on the new elements.
// It follows the following principle:
//   - If it's a whole new series, the title is the series' name. e.g. "Game of thrones"
//   - If there is (are) new season(s) (or new episodes in different seasons), the title is the series' name and the list of new Season number. e.g. "Game of thrones: Season 1-9"
//   - If there are only new episodes in a same season, the title is the series' name, season name and the list of new Season number. e.g. "Game of thrones: Season 1, Episodes 1-6 & 9"
//
// Finally the script aggregate consecutive number as much as possible, i.e. instead of having Episodes 1, 2, 3, 4 the script will squeeze in 1-4.
func buildNewSeriesItemFromSeriesNewItems(item jellyfin.NewlyAddedSeriesItem, app *app.ApplicationContext) string {
	if item.IsSeriesNew || len(item.NewSeasons) == 0 {
		return item.SeriesName
	}

	localizedSeason := app.Localizer.Localize("season", len(item.NewSeasons))

	newSeasonsNumber := []int{}
	newSeasonID := []string{}
	for seasonID, season := range item.NewSeasons {
		newSeasonsNumber = append(newSeasonsNumber, int(season.SeasonNumber))
		newSeasonID = append(newSeasonID, seasonID)
	}

	title := item.SeriesName + ": " + localizedSeason + " " + compressNumbers(newSeasonsNumber)

	if len(newSeasonID) == 1 && !item.NewSeasons[newSeasonID[0]].IsSeasonNew {
		newEpisodeNumber := []int{}
		for _, episode := range item.NewSeasons[newSeasonID[0]].Episodes {
			newEpisodeNumber = append(newEpisodeNumber, int(episode.EpisodeNumber))
		}
		localizedEpisode := app.Localizer.Localize("episode", len(newEpisodeNumber))
		title += " " + localizedEpisode + " " + compressNumbers(newEpisodeNumber)
	}

	return title
}

func buildNewMediaTemplateData(
	newMovies *[]jellyfin.MovieItem,
	newSeries *[]jellyfin.NewlyAddedSeriesItem,
	app *app.ApplicationContext) NewMediaTemplateData {
	newMoviesData := []NewMovieItemTemplateData{}
	newSeriesData := []NewSeriesItemTemplateData{}

	for _, newMovieItem := range *newMovies {
		newMoviesData = append(newMoviesData, NewMovieItemTemplateData{
			PosterURL:    newMovieItem.PosterURL,
			Name:         newMovieItem.Name,
			AdditionDate: newMovieItem.AdditionDate.Format("2006-01-02"),
			Overview:     newMovieItem.Overview,
			AddedOnLabel: app.Localizer.Localize("added_on", 1),
		})
	}

	for _, newSeriesItem := range *newSeries {
		newSeriesData = append(newSeriesData, NewSeriesItemTemplateData{
			PosterURL:      newSeriesItem.PosterURL,
			SeriesName:     newSeriesItem.SeriesName,
			AddedOnLabel:   app.Localizer.Localize("added_on", 1),
			AdditionDate:   newSeriesItem.AdditionDate.Format("2006-01-02"),
			Overview:       newSeriesItem.Overview,
			NewSeriesTitle: buildNewSeriesItemFromSeriesNewItems(newSeriesItem, app),
		})
	}
}

func BuildNewMediaEmailHTML(
	newMovies *[]jellyfin.MovieItem,
	newSeries *[]jellyfin.NewlyAddedSeriesItem,
	app *app.ApplicationContext,
) (string, error) {
	template, err := getNewMediaHTMLTemplate("new_media", app)

	if err != nil {
		// Error already logged
		return "", err
	}

}
