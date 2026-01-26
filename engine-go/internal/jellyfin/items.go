package jellyfin

import (
	"context"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

func (client *APIClient) GetRootFolderIdByName(folderName string, app *app.ApplicationContext) (string, error) {
	foldersItems, httpResponse, err := client.ItemsAPI.GetItems(context.Background()).
		Recursive(false).
		LocationTypes([]api.LocationType{api.LOCATIONTYPE_FILE_SYSTEM}).
		Execute()
	if err != nil {
		logHttpResponseError(httpResponse, err, app)
		return "", err
	}
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
