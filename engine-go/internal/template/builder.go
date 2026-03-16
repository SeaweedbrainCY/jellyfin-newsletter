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
	PosterURL            string
	Name                 string
	AddedOnLabel         string
	AdditionDate         string
	Overview             string
	IncludeItemOverviews bool
}
type newSeriesItemTemplateData struct {
	PosterURL            string
	SeriesName           string
	AddedOnLabel         string
	AdditionDate         string
	Overview             string
	NewSeriesTitle       string
	IncludeItemOverviews bool
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
	JellyfinOwnerName string
	UnsubscribeEmail  string
}

//go:embed themes/*/*.html
var templateHTMLThemesFS embed.FS

func CheckIfThemeIsAvailable(app *app.ApplicationContext) error {
	filePath := filepath.Join(
		"themes",
		"new_media",
		app.Config.EmailTemplate.Theme,
		app.Config.EmailTemplate.Theme+".html",
	)
	if _, err := template.ParseFS(templateHTMLThemesFS, filePath); err != nil {
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
	tmpl, err := template.New("Email").Option("missingkey=zero").ParseFS(themeFS, filePath)

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
		return strconv.Itoa(start)
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

	localizedSeason := app.Localizer.LocalizeWithPlural("season", len(item.NewSeasons))

	newSeasonsNumber := []int{}
	newSeasonID := []string{}
	for seasonID, season := range item.NewSeasons {
		newSeasonsNumber = append(newSeasonsNumber, int(season.SeasonNumber))
		newSeasonID = append(newSeasonID, seasonID)
	}

	title := item.SeriesName + ": " + localizedSeason + " " + compressNumbers(newSeasonsNumber)

	// Add episodes items
	if len(newSeasonID) == 1 && !item.NewSeasons[newSeasonID[0]].IsSeasonNew {
		newEpisodeNumber := []int{}
		for _, episode := range item.NewSeasons[newSeasonID[0]].Episodes {
			newEpisodeNumber = append(newEpisodeNumber, int(episode.EpisodeNumber))
		}
		localizedEpisode := app.Localizer.LocalizeWithPlural("episode", len(newEpisodeNumber))
		title += ", " + localizedEpisode + " " + compressNumbers(newEpisodeNumber)
	}

	return title
}

func buildStringTemplateWithPlaceholders(
	templateStr string,
	observedPeriodDays int,
	app *app.ApplicationContext,
) (string, error) {
	today := time.Now()
	todayDayNumber := int(today.Weekday())
	todayMonthNumber := int(today.Month())

	startDate := time.Now().AddDate(0, 0, observedPeriodDays*-1)
	startDayNumber := int(startDate.Weekday())
	startMonthNumber := int(startDate.Month())

	daysName := map[int]string{
		0: "sunday",
		1: "monday",
		2: "tuesday",
		3: "wednesday",
		4: "thursday",
		5: "friday",
		6: "saturday",
	}

	monthsName := map[int]string{
		1:  "january",
		2:  "february",
		3:  "march",
		4:  "april",
		5:  "may",
		6:  "june",
		7:  "july",
		8:  "august",
		9:  "september",
		10: "october",
		11: "november",
		12: "december",
	}

	placeholders := titlePlaceholders{
		Date:             today.Format("2006-01-02"),
		DayName:          app.Localizer.Localize(daysName[todayDayNumber]),
		DayNumber:        strconv.Itoa(todayDayNumber),
		MonthName:        app.Localizer.Localize(monthsName[todayMonthNumber]),
		MonthNumber:      strconv.Itoa(todayMonthNumber),
		Year:             today.Format("2006"),
		StartDate:        startDate.Format("2006-01-02"),
		StartDayName:     app.Localizer.Localize(daysName[startDayNumber]),
		StartDayNumber:   strconv.Itoa(startDayNumber),
		StartMonthName:   app.Localizer.Localize(monthsName[startMonthNumber]),
		StartMonthNumber: strconv.Itoa(todayMonthNumber),
		StartYear:        startDate.Format("2006"),
	}

	tmpl, err := template.New("title").Option("missingkey=zero").Parse(templateStr)
	if err != nil {
		app.Logger.Debug("Error while building title template", zap.String("templateStr", templateStr), zap.Error(err))
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, placeholders)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func buildFooterLabel(app *app.ApplicationContext) string {
	footerData := footerTemplateData{
		JellyfinOwnerName: app.Config.EmailTemplate.JellyfinOwnerName,
		UnsubscribeEmail:  app.Config.EmailTemplate.UnsubscribeEmail,
	}

	return app.Localizer.LocalizeWithTemplate("footer_label", footerData)
}

func sortJellyfinNewMovies(newJellyfinMovies *[]jellyfin.MovieItem, app *app.ApplicationContext) []jellyfin.MovieItem {
	newJellyfinMoviesSorted := slices.Clone(*newJellyfinMovies)
	slices.SortFunc(newJellyfinMoviesSorted, func(a, b jellyfin.MovieItem) int {
		switch app.Config.EmailTemplate.SortMode {
		case "name_asc":
			return strings.Compare(a.Name, b.Name)
		case "name_desc":
			return strings.Compare(b.Name, a.Name)
		case "date_desc":
			return b.AdditionDate.Compare(*a.AdditionDate)
		// date_asc is the default option
		default:
			return a.AdditionDate.Compare(*b.AdditionDate)
		}
	})
	return newJellyfinMoviesSorted
}

func sortJellyfinNewSeriesItems(
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem,
	app *app.ApplicationContext,
) []jellyfin.NewlyAddedSeriesItem {
	newJellyfinSeriesSorted := slices.Clone(*newJellyfinSeries)
	slices.SortFunc(newJellyfinSeriesSorted, func(a, b jellyfin.NewlyAddedSeriesItem) int {
		switch app.Config.EmailTemplate.SortMode {
		case "name_asc":
			return strings.Compare(a.SeriesName, b.SeriesName)
		case "name_desc":
			return strings.Compare(b.SeriesName, a.SeriesName)
		case "date_desc":
			return getAdditionDateForSeries(b).Compare(getAdditionDateForSeries(a))
		// date_asc is the default option
		default:
			return getAdditionDateForSeries(a).Compare(getAdditionDateForSeries(b))
		}
	})
	return newJellyfinSeriesSorted
}

func shouldOverviewsBeDisplayed(itemsCount int, app *app.ApplicationContext) bool {
	switch app.Config.EmailTemplate.DisplayOverviewMaxItems {
	case -1:
		return false
	case 0:
		return true
	default:
		return itemsCount < app.Config.EmailTemplate.DisplayOverviewMaxItems
	}
}

// The addition date is not necessarily the series's addition date.
// If the new element is a season or an episode, it's its addition date that we want
// If there are several new elements, the newest date is kept.
func getAdditionDateForSeries(seriesItem jellyfin.NewlyAddedSeriesItem) time.Time {
	if seriesItem.IsSeriesNew {
		return seriesItem.AdditionDate
	}
	additionDate := time.Date(1970, 01, 01, 00, 00, 0, 0, time.UTC)
	for _, season := range seriesItem.NewSeasons {
		if season.IsSeasonNew {
			if season.AdditionDate.After(additionDate) {
				additionDate = season.AdditionDate
			}
			continue
		}
		// Not a new season
		for _, episode := range season.Episodes {
			if episode.AdditionDate.After(additionDate) {
				additionDate = episode.AdditionDate
			}
		}
	}
	return additionDate
}

func buildNewMediaTemplateData(
	newJellyfinMovies *[]jellyfin.MovieItem,
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem,
	movieCount int32,
	episodesCount int32,
	app *app.ApplicationContext) (*newMediaTemplateData, error) {
	newMoviesData := []newMovieItemTemplateData{}
	newSeriesData := []newSeriesItemTemplateData{}

	htmlDir := "ltr"
	if slices.Contains(
		[]string{"ar", "he", "fa", "ur", "ku", "ps", "yi", "dv", "qrc"},
		app.Config.EmailTemplate.Language,
	) {
		htmlDir = "rtl"
	}

	newJellyfinMoviesSorted := sortJellyfinNewMovies(newJellyfinMovies, app)

	newJellyfinSeriesSorted := sortJellyfinNewSeriesItems(newJellyfinSeries, app)

	displayMovieOverviews := shouldOverviewsBeDisplayed(len(newJellyfinMoviesSorted), app)
	displaySeriesOverviews := shouldOverviewsBeDisplayed(len(newJellyfinSeriesSorted), app)

	for _, newMovieItem := range newJellyfinMoviesSorted {
		newMoviesData = append(newMoviesData, newMovieItemTemplateData{
			PosterURL:            newMovieItem.PosterURL,
			Name:                 newMovieItem.Name,
			AdditionDate:         newMovieItem.AdditionDate.Format("2006-01-02"),
			Overview:             newMovieItem.Overview,
			AddedOnLabel:         app.Localizer.Localize("added_on"),
			IncludeItemOverviews: displayMovieOverviews,
		})
	}

	for _, newSeriesItem := range newJellyfinSeriesSorted {
		newSeriesData = append(newSeriesData, newSeriesItemTemplateData{
			PosterURL:            newSeriesItem.PosterURL,
			SeriesName:           newSeriesItem.SeriesName,
			AddedOnLabel:         app.Localizer.Localize("added_on"),
			AdditionDate:         getAdditionDateForSeries(newSeriesItem).Format("2006-01-02"),
			Overview:             newSeriesItem.Overview,
			NewSeriesTitle:       buildNewSeriesItemFromSeriesNewItems(newSeriesItem, app),
			IncludeItemOverviews: displaySeriesOverviews,
		})
	}

	title, err := buildStringTemplateWithPlaceholders(
		app.Config.EmailTemplate.Title,
		app.Config.Jellyfin.ObservedPeriodDays,
		app,
	)
	if err != nil {
		return nil, err
	}

	subtitle, err := buildStringTemplateWithPlaceholders(
		app.Config.EmailTemplate.Subtitle,
		app.Config.Jellyfin.ObservedPeriodDays,
		app,
	)
	if err != nil {
		return nil, err
	}

	data := newMediaTemplateData{
		HTMLLang:                     app.Config.EmailTemplate.Language,
		HTMLDir:                      htmlDir,
		Title:                        title,
		Subtitle:                     subtitle,
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
		MoviesLabel:                  app.Localizer.LocalizeWithPlural("movies", int(movieCount)),
		SeriesLabel:                  app.Localizer.LocalizeWithPlural("episode", int(episodesCount)),
		FooterLabel:                  buildFooterLabel(app),
		FooterProjectLinkLabel:       "Jellyfin Newsletter",
		FooterOpenSourceProjectLabel: app.Localizer.Localize("footer_project_open_source"),
		FooterDevelopedByLabel:       app.Localizer.Localize("footer_developed_by"),
		FooterLicenceAndCopyright:    app.Localizer.Localize("license_and_copyright"),
	}
	return &data, nil
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

	tmplData, err := buildNewMediaTemplateData(newMovies, newSeries, movieCount, episodesCount, app)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplData)

	if err != nil {
		app.Logger.Error("An error occurred while populating the email HTML template", zap.Error(err))
		return "", err
	}

	return buf.String(), nil
}
