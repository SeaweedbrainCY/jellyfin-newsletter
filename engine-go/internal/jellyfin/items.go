package jellyfin

import (
	"context"
	"strconv"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

type Nullable[T any] interface {
	IsSet() bool
	Get() *T
}

func (client *APIClient) GetRootFolderIDByName(folderName string, app *app.ApplicationContext) (string, error) {
	foldersItems, httpResponse, err := client.ItemsAPI.GetItems(context.Background()).
		Recursive(false).
		LocationTypes([]api.LocationType{api.LOCATIONTYPE_FILE_SYSTEM}).
		Execute()
	if err != nil {
		logHTTPResponseError(httpResponse, err, app)
		return "", err
	}
	defer httpResponse.Body.Close()
	if !foldersItems.HasItems() {
		app.Logger.Warn(
			"No folders found. This could happen if Jellyfin has no collection or folder at all. Media should belong in folders but none are found.",
		)
		return "", ErrItemsNotFound
	}
	for _, item := range foldersItems.GetItems() {
		if *item.Name.Get() == folderName {
			return *item.Id, nil
		}
	}
	app.Logger.Warn("Folder not found. Will ignore it.", zap.String("folderName", folderName))
	return "", ErrItemsNotFound
}

func OrDefault[T any](n Nullable[T], def T) T {
	if n.IsSet() {
		return *n.Get()
	}
	return def
}

