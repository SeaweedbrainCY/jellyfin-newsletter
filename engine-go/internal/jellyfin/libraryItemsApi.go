package jellyfin

import (
	"context"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

type libraryItemAPI struct {
	jellyfinAPI.LibraryAPI
}

type Nullable[T any] interface {
	IsSet() bool
	Get() *T
}

func (libraryAPI libraryItemAPI) GetItemsStats(app *app.ApplicationContext) (int32, int32, error) {
	itemsCounts, httpResponse, httpErr := libraryAPI.GetItemCounts(context.Background()).Execute()

	err := checkHTTPRequest("GetItemsStats", httpResponse, httpErr, app.Logger)
	if err != nil {
		return 0, 0, err
	}

	defer httpResponse.Body.Close()
	defer httpResponse.Body.Close()
	return *itemsCounts.MovieCount, *itemsCounts.EpisodeCount, nil
}

func (libraryAPI libraryItemAPI) GetMoviesItemsByFolderID(
	folderID string,
	app *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	movies, getMoviesHTTPResponse, httpErr := libraryAPI.GetItems(context.Background()).
		Recursive(false).
		EnableImages(false).
		ParentId(folderID).
		LocationTypes([]jellyfinAPI.LocationType{jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM}).
		IsMovie(true).
		Fields([]jellyfinAPI.ItemFields{"DateCreated", "ProviderIds", "Id", "Name", "ProductionYear"}).
		Execute()

	err := checkHTTPRequest("GetMoviesItemsByFolderID", getMoviesHTTPResponse, httpErr, app.Logger)
	if err != nil {
		return nil, err
	}
	defer getMoviesHTTPResponse.Body.Close()
	app.Logger.Debug("Successfully retrieved all movies from API", zap.String("Folder", folderID), zap.Any("Movies", movies.Items))
	return &movies.Items, nil
}

func (libraryAPI libraryItemAPI) GetAllItemsByFolderID(
	folderID string,
	app *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	items, httpResponse, httpErr := libraryAPI.GetItems(context.Background()).
		Recursive(true).
		EnableImages(false).
		ParentId(folderID).
		Fields([]jellyfinAPI.ItemFields{"DateCreated", "ProviderIds", "Id", "Name", "ProductionYear", "IndexNumber", "SeriesId", "Type", "SeasonId"}).
		Execute()

	err := checkHTTPRequest("GetAllItemsByFolderID", httpResponse, httpErr, app.Logger)
	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()
	app.Logger.Debug("Successfully retrieved all items from API", zap.String("Folder", folderID), zap.Any("Items", items.Items))
	return &items.Items, nil
}

func (libraryAPI libraryItemAPI) GetRootFolderIDByName(folderName string, app *app.ApplicationContext) (string, error) {
	foldersItems, httpResponse, httpErr := libraryAPI.GetItems(context.Background()).
		Recursive(false).
		LocationTypes([]jellyfinAPI.LocationType{jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM}).
		Execute()

	err := checkHTTPRequest("GetRootFolderIDByName", httpResponse, httpErr, app.Logger)
	if err != nil {
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
		if value == nil {
			return ""
		}
		return *value
	}
	return ""
}
