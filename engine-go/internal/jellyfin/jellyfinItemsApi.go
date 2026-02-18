package jellyfin

import (
	"context"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

type jellyfinItemsAPI struct {
	jellyfinAPI.ItemsAPI
}

type Nullable[T any] interface {
	IsSet() bool
	Get() *T
}

func (itemsAPI jellyfinItemsAPI) GetMoviesItemsByFolderID(
	folderID string,
	recursive bool,
	app *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	movies, getMoviesHTTPResponse, err := itemsAPI.GetItems(context.Background()).
		Recursive(recursive).
		ParentId(folderID).
		LocationTypes([]jellyfinAPI.LocationType{jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM}).
		IsMovie(true).
		Fields([]jellyfinAPI.ItemFields{"DateCreated", "ProviderIds", "Id", "Name", "ProductionYear"}).
		Execute()

	if err != nil {
		logHTTPResponseError(getMoviesHTTPResponse, err, app)
		return nil, err
	}
	defer getMoviesHTTPResponse.Body.Close()
	return &movies.Items, nil
}

func (itemsAPI jellyfinItemsAPI) GetAllItemsByFolderID(
	folderID string,
	app *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	items, httpResponse, err := itemsAPI.GetItems(context.Background()).
		Recursive(true).
		ParentId(folderID).
		Fields([]jellyfinAPI.ItemFields{"DateCreated", "ProviderIds", "Id", "Name", "ProductionYear", "IndexNumber", "SeriesId", "Type", "SeasonId"}).
		Execute()
	if err != nil {
		logHTTPResponseError(httpResponse, err, app)
		return nil, err
	}
	return &items.Items, nil
}

func (itemsAPI jellyfinItemsAPI) GetRootFolderIDByName(folderName string, app *app.ApplicationContext) (string, error) {
	foldersItems, httpResponse, err := itemsAPI.GetItems(context.Background()).
		Recursive(false).
		LocationTypes([]jellyfinAPI.LocationType{jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM}).
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
	if n.IsSet() && n.Get() != nil {
		return *n.Get()
	}
	return def
}

func getTMDBIDIfExist(item *jellyfinAPI.BaseItemDto) string {
	if value, ok := item.ProviderIds["Tmdb"]; ok {
		return value
	}
	return ""
}
