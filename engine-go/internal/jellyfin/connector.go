package jellyfin

import (
	"net/http"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"go.uber.org/zap"
)

func (client *APIClient) TestConnection(app *app.ApplicationContext) error {
	pingAnswer, pingHTTPReponse, err := client.SystemAPI.PingSystem()
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

	systemInfo, httpStatusCode, systemInfoerr := client.SystemAPI.GetSystemInformation()
	if systemInfoerr != nil {
		return systemInfoerr
	}

	app.Logger.Info(
		"Successfully connected to Jellyfin",
		zap.Int("http_status", httpStatusCode),
		zap.String("apiVersion", systemInfo.APIVersion),
		zap.String("serverName", systemInfo.ServerName),
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
