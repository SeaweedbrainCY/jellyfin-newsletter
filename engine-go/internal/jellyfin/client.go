package jellyfin

import (
	"net/http"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type SystemAPIInterface interface {
	PingSystem() (string, *http.Response, error)
	GetSystemInformation() (*SystemInfo, int, error)
}

type ItemsAPIInterface interface {
	GetMoviesItemsByFolderID(
		folderID string,
		recursive bool,
		app *app.ApplicationContext,
	) (*[]jellyfinAPI.BaseItemDto, error)
	GetRootFolderIDByName(folderName string, app *app.ApplicationContext) (string, error)
	GetAllItemsByFolderID(
		folderID string,
		app *app.ApplicationContext,
	) (*[]jellyfinAPI.BaseItemDto, error)
}

type APIClient struct {
	SystemAPI SystemAPIInterface
	ItemsAPI  ItemsAPIInterface
}

func NewJellyfinAPIClient(app *app.ApplicationContext) *APIClient {
	headerToken := "MediaBrowser Token=\"" + app.Config.Jellyfin.APIKey + "\""
	config := &jellyfinAPI.Configuration{
		Servers:       jellyfinAPI.ServerConfigurations{{URL: app.Config.Jellyfin.URL}},
		DefaultHeader: map[string]string{"Authorization": headerToken},
	}
	client := jellyfinAPI.NewAPIClient(config)
	return &APIClient{
		SystemAPI: JellyfinSystemAPI{
			client.SystemAPI,
		},
		ItemsAPI: jellyfinItemsAPI{
			client.ItemsAPI,
		},
	}
}
