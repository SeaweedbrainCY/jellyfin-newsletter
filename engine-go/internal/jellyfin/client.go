package jellyfin

import (
	"context"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type ItemsAPIInterface interface {
	GetItems(ctx context.Context) jellyfinAPI.ApiGetItemsRequest
}

type SystemAPIInterface interface {
	GetSystemInfo(ctx context.Context) jellyfinAPI.ApiGetSystemInfoRequest
	PostPingSystem(ctx context.Context) jellyfinAPI.ApiPostPingSystemRequest
}

type APIClient struct {
	ItemsAPI  ItemsAPIInterface
	SystemAPI SystemAPIInterface
}

func GetJellyfinAPIClient(app *app.ApplicationContext) *APIClient {
	headerToken := "MediaBrowser Token=\"" + app.Config.Jellyfin.APIKey + "\""
	config := &jellyfinAPI.Configuration{
		Servers:       jellyfinAPI.ServerConfigurations{{URL: app.Config.Jellyfin.URL}},
		DefaultHeader: map[string]string{"Authorization": headerToken},
	}
	client := jellyfinAPI.NewAPIClient(config)
	return &APIClient{
		ItemsAPI: client.ItemsAPI,
	}
}
