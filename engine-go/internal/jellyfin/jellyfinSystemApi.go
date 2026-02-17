package jellyfin

import (
	"context"
	"net/http"

	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type jellyfinSystemAPI struct {
	jellyfinAPI.SystemAPI
}

type SystemInfo struct {
	APIVersion string
	ServerName string
}

func (systemAPI jellyfinSystemAPI) PingSystem() (string, *http.Response, error) {
	return systemAPI.PostPingSystem(context.Background()).Execute()
}

func (systemAPI jellyfinSystemAPI) GetSystemInformation() (*SystemInfo, int, error) {
	systemInfo, systemInfoHTTPReponse, err := systemAPI.GetSystemInfo(context.Background()).Execute()
	statusCode := 0
	if systemInfoHTTPReponse != nil {
		statusCode = systemInfoHTTPReponse.StatusCode
		defer systemInfoHTTPReponse.Body.Close()
	}
	if err != nil {
		return nil, statusCode, err
	}

	return &SystemInfo{
		APIVersion: OrDefault(systemInfo.Version, "Unknown"),
		ServerName: OrDefault(systemInfo.ServerName, "Unknown"),
	}, statusCode, nil
}
