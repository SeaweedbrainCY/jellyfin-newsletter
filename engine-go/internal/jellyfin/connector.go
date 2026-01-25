package jellyfin

import (
	"context"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

func (client *APIClient) TestConnection(app *app.ApplicationContext) error {
	pingAnswer, r, err := client.SystemAPI.PostPingSystem(context.Background()).Execute()
	if err != nil {
		app.Logger.Error(
			"Ping to Jellyfin API failed. Check for network issues.",
			zap.Int("ping_HTTP_status", r.StatusCode),
			zap.String("response", pingAnswer),
		)
		return err
	}
	app.Logger.Debug(
		"Successfully pinged the Jellyfin API",
		zap.Int("ping_HTTP_status", r.StatusCode),
		zap.String("response", pingAnswer),
	)

	systemInfo, r, err := client.SystemAPI.GetSystemInfo(context.Background()).Execute()
	if err != nil {
		app.Logger.Error(
			"Failed to connect to the Jellyfin API",
			zap.Int("http_status", r.StatusCode),
		)
		return err
	}

	apiVersion := "Unknown"
	serverName := "Unknown"
	if api.NullableString.IsSet(systemInfo.Version) {
		apiVersion = *api.NullableString.Get(systemInfo.Version)
	}
	if systemInfo.ServerName.IsSet() {
		serverName = *systemInfo.ServerName.Get()
	}

	app.Logger.Info(
		"Successfully connected to Jellyfin",
		zap.Int("http_status", r.StatusCode),
		zap.String("apiVersion", apiVersion),
		zap.String("serverName", serverName),
	)

	return err
}
