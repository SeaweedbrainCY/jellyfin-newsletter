package jellyfin

import (
	"context"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type libraryItemAPI struct {
	jellyfinAPI.LibraryAPI
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
