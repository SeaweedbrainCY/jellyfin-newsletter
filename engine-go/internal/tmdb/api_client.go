package tmdb

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"go.uber.org/zap"
)

type APIInterface interface {
	GetMediaByID(id string, mediaType MediaType) (*GetMediaHTTPResponse, error)
	GetMediaByTitle(mediaName string, mediaType MediaType)
}

type APIClient struct {
	APIKey string
	Lang   string
	Logger *zap.Logger
}

func InitTMDBApiClient(app *app.ApplicationContext) APIClient {
	return APIClient{
		APIKey: app.Config.TMDB.APIKey,
		Lang:   app.Config.EmailTemplate.Language,
		Logger: app.Logger,
	}
}
