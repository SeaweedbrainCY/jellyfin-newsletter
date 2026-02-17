package jellyfin

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type MockJellyfinItemsAPI struct {
	ExecuteGetMoviesItemsByFolderID func() (*[]jellyfinAPI.BaseItemDto, error)
	ExecuteGetRootFolderIDByName    func() (string, error)
	ExecuteGetAllItemsByFolderID    func() (*[]jellyfinAPI.BaseItemDto, error)
}

func (m MockJellyfinItemsAPI) GetMoviesItemsByFolderID(
	_ string,
	_ bool,
	_ *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return m.ExecuteGetMoviesItemsByFolderID()
}

func (m MockJellyfinItemsAPI) GetAllItemsByFolderID(
	_ string,
	_ *app.ApplicationContext,
) (*[]jellyfinAPI.BaseItemDto, error) {
	return m.ExecuteGetAllItemsByFolderID()
}

func (m MockJellyfinItemsAPI) GetRootFolderIDByName(_ string, _ *app.ApplicationContext) (string, error) {
	return m.ExecuteGetRootFolderIDByName()
}

func Ptr[T any](v T) *T {
	return &v
}
