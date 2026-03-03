package template

import (
	"os"
	"path/filepath"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"go.uber.org/zap"
)

type NewMovieItemTemplateData struct {
	PosterURL    string
	Name         string
	AddedOnLabel string
	AdditionDate string
	Overview     string
}
type NewSeriesItemTemplateData struct {
	PosterURL      string
	SeriesName     string
	AddedOnLabel   string
	AdditionDate   string
	Overview       string
	NewSeriesTitle string
}

type NewMediaTemplateData struct {
	HTMLLang                     string
	HTMLDir                      string
	Title                        string
	Subtitle                     string
	JellyfinURL                  string
	DiscoverNowLabel             string
	DisplayNewMovies             bool
	NewFilmLabel                 string
	NewMovies                    []NewMovieItemTemplateData
	DisplayNewSeries             bool
	NewSeriesLabel               string
	NewSeries                    []NewSeriesItemTemplateData
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

func getNewMediaHTMLTemplate(themes_type string, app *app.ApplicationContext) (string, error) {
	filePath := filepath.Join("themes", themes_type, app.Config.EmailTemplate.Theme)
	template, err := os.ReadFile(filePath)

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

func buildNewMediaTemplateData(
	newMovies *[]jellyfin.MovieItem,
	newSeries *[]jellyfin.NewlyAddedSeriesItem,
	app *app.ApplicationContext) NewMediaTemplateData {
	newMoviesData := []NewMovieItemTemplateData{}

	for _, newMovieItem := range *newMovies {
		newMoviesData = append(newMoviesData, NewMovieItemTemplateData{
			PosterURL:    newMovieItem.PosterURL,
			Name:         newMovieItem.Name,
			AdditionDate: newMovieItem.AdditionDate.Format("2006-01-02"),
			Overview:     newMovieItem.Overview,
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
