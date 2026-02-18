package tmdb

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"go.uber.org/zap"
)

type TMDBAPIInterface interface {
	GetMediaByID(id string, mediaType MediaType) (*GetMediaHTTPResponse, error)
	GetMediaByTitle(mediaName string, mediaType MediaType)
}

type TMDBAPIClient struct {
	APIKey string
	Lang   string
	Logger *zap.Logger
}

func InitTMDBApiClient(app *app.ApplicationContext) TMDBAPIClient {
	return TMDBAPIClient{
		APIKey: app.Config.TMDB.APIKey,
		Lang:   app.Config.EmailTemplate.Language,
		Logger: app.Logger,
	}
}
