package jellyfin

import (
	"context"
	"net/http"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

func (client *APIClient) TestConnection(app *app.ApplicationContext) error {
	pingAnswer, pingHTTPReponse, err := client.SystemAPI.PostPingSystem(context.Background()).Execute()
	statusCode := 0
	if pingHTTPReponse != nil {
		statusCode = pingHTTPReponse.StatusCode
		defer pingHTTPReponse.Body.Close()
	}
	if err != nil {
		app.Logger.Error(
			"Ping to Jellyfin API failed. Check for network issues.",
			zap.Int("ping_HTTP_status", statusCode),
			zap.String("response", pingAnswer),
			zap.Error(err),
		)
		return err
	}
	app.Logger.Debug(
		"Successfully pinged the Jellyfin API",
		zap.Int("ping_HTTP_status", statusCode),
		zap.String("response", pingAnswer),
	)

	systemInfo, systemInfoHTTPReponse, err := client.SystemAPI.GetSystemInfo(context.Background()).Execute()

	statusCode = 0
	if pingHTTPReponse != nil {
		statusCode = systemInfoHTTPReponse.StatusCode
		defer systemInfoHTTPReponse.Body.Close()
	}
	if err != nil {
		app.Logger.Error(
			"Failed to connect to the Jellyfin API",
			zap.Int("http_status", statusCode),
			zap.Error(err),
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
		zap.Int("http_status", statusCode),
		zap.String("apiVersion", apiVersion),
		zap.String("serverName", serverName),
	)

	return err
}

func logHTTPResponseError(httpResponse *http.Response, err error, app *app.ApplicationContext) {
	statusCode := 0
	if httpResponse != nil {
		statusCode = httpResponse.StatusCode
		defer httpResponse.Body.Close()
	}
	app.Logger.Error("Getting root Items failed.", zap.Int("httpStatusCode", statusCode), zap.Error(err))
}
