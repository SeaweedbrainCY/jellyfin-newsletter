package jellyfin

import (
	"context"
	"testing"

	jellyfinAPI "github.com/sj14/jellyfin-go/api"
)

type MockItemsAPI struct {
	GetItemsFunc func(ctx context.Context) jellyfinAPI.ApiGetItemsRequest
}

type MockApiGetItemsRequest struct {
	ExecuteFunc func() jellyfinAPI.BaseItemDtoQueryResult
}

func (m *MockItemsAPI) GetItems(ctx context.Context) jellyfinAPI.ApiGetItemsRequest {
	return m.GetItemsFunc(ctx)
}

func (m *MockApiGetItemsRequest) Execute() jellyfinAPI.BaseItemDtoQueryResult {
	return m.ExecuteFunc()
}

func TestGetRecentlyAddedMoviesByFolder(t *testing.T) {

}
