package jellyfin

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type APIClient struct {
	*jellyfinAPI.APIClient
}

func GetJellyfinAPIClient(app *app.ApplicationContext) *APIClient {
	headerToken := "`MediaBrowser Token=\"" + app.Config.Jellyfin.APIKey + "\"`"
	config := &jellyfinAPI.Configuration{
		Servers:       jellyfinAPI.ServerConfigurations{{URL: app.Config.Jellyfin.URL}},
		DefaultHeader: map[string]string{"Authorization": headerToken},
	}

	return &APIClient{jellyfinAPI.NewAPIClient(config)}
}
