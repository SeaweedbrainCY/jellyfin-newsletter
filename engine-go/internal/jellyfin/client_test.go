package jellyfin

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type MockJellyfinLibraryAPI struct {
	ExecuteGetMoviesItemsByFolderID func() (*[]jellyfinAPI.BaseItemDto, error)
	ExecuteGetRootFolderIDByName    func() (string, error)
	ExecuteGetAllItemsByFolderID    func() (*[]jellyfinAPI.BaseItemDto, error)
	ExecuteGetItemsStats            func() (int32, int32, error)
}

func (m MockJellyfinLibraryAPI) GetMoviesItemsByFolderID(
	_ string,
	_ *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return m.ExecuteGetMoviesItemsByFolderID()
}

func (m MockJellyfinLibraryAPI) GetAllItemsByFolderID(
	_ string,
	_ *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return m.ExecuteGetAllItemsByFolderID()
}

func (m MockJellyfinLibraryAPI) GetRootFolderIDByName(_ string, _ *app.ApplicationContext) (string, error) {
	return m.ExecuteGetRootFolderIDByName()
}

func (m MockJellyfinLibraryAPI) GetItemsStats(_ *app.ApplicationContext) (int32, int32, error) {
	return m.ExecuteGetItemsStats()
}
