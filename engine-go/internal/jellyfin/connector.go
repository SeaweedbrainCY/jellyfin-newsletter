package jellyfin

import (
	"errors"
	"io"
	"net/http"
	"strconv"

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

func checkHTTPRequest(ctx string, resp *http.Response, httpErr error, logger *zap.Logger) error {
	if httpErr != nil {
		logger.Error(
			"HTTP request failed",
			zap.String("context", ctx),
			zap.Error(httpErr),
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, bodyErr := io.ReadAll(resp.Body)

		fields := []zap.Field{
			zap.Int("status", resp.StatusCode),
			zap.String("context", ctx),
		}

		if bodyErr == nil {
			fields = append(fields, zap.String("body", string(body)))
		} else {
			fields = append(fields, zap.NamedError("body_read_error", bodyErr))
		}

		logger.Error("Unexpected HTTP response", fields...)
		return errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
