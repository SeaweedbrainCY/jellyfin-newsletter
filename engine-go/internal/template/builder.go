package template

import (
	"bytes"
	"embed"
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

type footerTemplateData struct {
	ownerName        string
	unsubscribeEmail string
}

func CheckIfThemeIsAvailable(themeFS embed.FS, app *app.ApplicationContext) error {
	filePath := filepath.Join(
		"themes",
		"new_media",
		app.Config.EmailTemplate.Theme,
		app.Config.EmailTemplate.Theme+".html",
	)
	if _, err := template.ParseFS(themeFS, filePath); err != nil {
		return err
	}
	return nil
}

func getNewMediaHTMLTemplate(
	themesType string,
	themeFS embed.FS,
	app *app.ApplicationContext,
) (*template.Template, error) {
	filePath := filepath.Join(
		"themes",
		themesType,
		app.Config.EmailTemplate.Theme,
		app.Config.EmailTemplate.Theme+".html",
	)
	tmpl, err := template.ParseFS(themeFS, filePath)

	if err != nil {
		app.Logger.Error(
			"Impossible to open the template file. Email HTML building will fail.",
			zap.String("filePath", filePath),
			zap.Error(err),
		)
		return nil, err
	}

	return tmpl, nil
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

func buildStringTemplateWithPlaceholders(titleTemplate string, observedPeriodDays int) string {
	today := time.Now()
	todayDayNumber := today.Format("2")
	todayMonthNumber := today.Format("1")

	startDate := time.Now().Add(time.Duration(observedPeriodDays) * -1)
	startDayNumber := startDate.Format("2")
	startMonthNumber := startDate.Format("1")

	daysName := map[string]string{
		"1": "monday",
		"2": "tuesday",
		"3": "wednesday",
		"4": "thursday",
		"5": "friday",
		"6": "saturday",
		"7": "sunday",
	}

	monthsName := map[string]string{
		"1":  "january",
		"2":  "february",
		"3":  "march",
		"4":  "april",
		"5":  "may",
		"6":  "june",
		"7":  "july",
		"8":  "august",
		"9":  "september",
		"10": "october",
		"11": "november",
		"12": "december",
	}

	placeholders := titlePlaceholders{
		Date:             today.Format("2006-01-02"),
		DayName:          daysName[todayDayNumber],
		DayNumber:        todayDayNumber,
		MonthName:        monthsName[todayMonthNumber],
		MonthNumber:      todayMonthNumber,
		Year:             today.Format("2006"),
		StartDate:        startDate.Format("2006-01-02"),
		StartDayName:     daysName[startDayNumber],
		StartDayNumber:   todayDayNumber,
		StartMonthName:   monthsName[startMonthNumber],
		StartMonthNumber: todayMonthNumber,
		StartYear:        startDate.Format("2006"),
	}

	tmpl := template.Must(template.New("title").Parse(titleTemplate))
	var buf bytes.Buffer
	tmpl.Execute(&buf, placeholders)

	return buf.String()
}

func buildFooterLabel(label string, app *app.ApplicationContext) string {
	footerData := footerTemplateData{
		ownerName:        app.Config.EmailTemplate.JellyfinOwnerName,
		unsubscribeEmail: app.Config.EmailTemplate.UnsubscribeEmail,
	}
	tmpl := template.Must(template.New("footer").Parse(label))
	var buf bytes.Buffer
	tmpl.Execute(&buf, footerData)

	return buf.String()
}

func buildNewMediaTemplateData(
	newMovies *[]jellyfin.MovieItem,
	newSeries *[]jellyfin.NewlyAddedSeriesItem,
	movieCount int32,
	episodesCount int32,
	app *app.ApplicationContext) newMediaTemplateData {
	newMoviesData := []newMovieItemTemplateData{}
	newSeriesData := []newSeriesItemTemplateData{}

	HTMLdir := "ltr"
	if slices.Contains(
		[]string{"ar", "he", "fa", "ur", "ku", "ps", "yi", "dv", "qrc"},
		app.Config.EmailTemplate.Language,
	) {
		HTMLdir = "rtl"
	}

	for _, newMovieItem := range *newMovies {
		newMoviesData = append(newMoviesData, newMovieItemTemplateData{
			PosterURL:    newMovieItem.PosterURL,
			Name:         newMovieItem.Name,
			AdditionDate: newMovieItem.AdditionDate.Format("2006-01-02"),
			Overview:     newMovieItem.Overview,
			AddedOnLabel: app.Localizer.Localize("added_on"),
		})
	}

	for _, newSeriesItem := range *newSeries {
		newSeriesData = append(newSeriesData, newSeriesItemTemplateData{
			PosterURL:      newSeriesItem.PosterURL,
			SeriesName:     newSeriesItem.SeriesName,
			AddedOnLabel:   app.Localizer.Localize("added_on"),
			AdditionDate:   newSeriesItem.AdditionDate.Format("2006-01-02"),
			Overview:       newSeriesItem.Overview,
			NewSeriesTitle: buildNewSeriesItemFromSeriesNewItems(newSeriesItem, app),
		})
	}

	return newMediaTemplateData{
		HTMLLang: app.Config.EmailTemplate.Language,
		HTMLDir:  HTMLdir,
		Title: buildStringTemplateWithPlaceholders(
			app.Config.EmailTemplate.Title,
			app.Config.Jellyfin.ObservedPeriodDays,
		),
		Subtitle: buildStringTemplateWithPlaceholders(
			app.Config.EmailTemplate.Subtitle,
			app.Config.Jellyfin.ObservedPeriodDays,
		),
		JellyfinURL:                  app.Config.EmailTemplate.JellyfinURL,
		DiscoverNowLabel:             app.Localizer.Localize("discover_now"),
		DisplayNewMovies:             len(newMoviesData) > 0,
		NewFilmLabel:                 app.Localizer.Localize("new_film"),
		NewMovies:                    newMoviesData,
		DisplayNewSeries:             len(newSeriesData) > 0,
		NewSeriesLabel:               app.Localizer.Localize("new_tvs"),
		NewSeries:                    newSeriesData,
		CurrentlyAvailableLabel:      app.Localizer.Localize("currently_available"),
		MoviesCount:                  strconv.Itoa(int(movieCount)),
		SeriesCount:                  strconv.Itoa(int(episodesCount)),
		MoviesLabel:                  app.Localizer.Localize("movies", int(movieCount)),
		SeriesLabel:                  app.Localizer.Localize("episode", int(episodesCount)),
		FooterLabel:                  buildFooterLabel(app.Localizer.Localize("footer_label"), app),
		FooterProjectLinkLabel:       "Jellyfin Newsletter",
		FooterOpenSourceProjectLabel: app.Localizer.Localize("footer_project_open_source"),
		FooterDevelopedByLabel:       app.Localizer.Localize("footer_developed_by"),
		FooterLicenceAndCopyright:    app.Localizer.Localize("license_and_copyright"),
	}
}

func BuildNewMediaEmailHTML(
	newMovies *[]jellyfin.MovieItem,
	newSeries *[]jellyfin.NewlyAddedSeriesItem,
	movieCount int32,
	episodesCount int32,
	app *app.ApplicationContext,
	themeFS embed.FS,
) (string, error) {
	tmpl, err := getNewMediaHTMLTemplate("new_media", themeFS, app)

	if err != nil {
		// Error already logged
		return "", err
	}

	tmplData := buildNewMediaTemplateData(newMovies, newSeries, movieCount, episodesCount, app)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)

	if err != nil {
		app.Logger.Error("An error occurred while populating the email HTML template", zap.Error(err))
		return "", err
	}

	return buf.String(), nil
}
