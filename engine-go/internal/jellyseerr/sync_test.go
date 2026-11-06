package jellyseerr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchRecipients_Disabled(t *testing.T) {
	recipients, err := FetchRecipients(context.Background(), nil, nil)
	require.NoError(t, err)
	assert.Nil(t, recipients)

	cfg := &config.JellyseerrConfig{Enabled: false}
	recipients, err = FetchRecipients(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.Nil(t, recipients)
}

func TestFetchRecipients_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-api-key", r.Header.Get("X-Api-Key"))
		assert.Equal(t, "/api/v1/user", r.URL.Path)

		resp := UsersResponse{
			Results: []User{
				{
					ID:               1,
					Email:            "fabian@example.com",
					DisplayName:      "Fabian",
					JellyfinUsername: "fabian",
				},
				{
					ID:               2,
					Email:            "excluded@example.com",
					DisplayName:      "Werner",
					JellyfinUsername: "werner",
				},
				{
					ID:               3,
					Email:            "", // no email, should be skipped
					DisplayName:      "NoEmailUser",
					JellyfinUsername: "noemail",
				},
				{
					ID:               4,
					Email:            "fabian@example.com", // duplicate email, should be skipped
					DisplayName:      "Fabian Duplicate",
					JellyfinUsername: "fabian2",
				},
				{
					ID:               5,
					Email:            "lorenz@example.com",
					DisplayName:      "Lorenz",
					JellyfinUsername: "lorenz",
				},
			},
		}
		resp.PageInfo.Pages = 1
		resp.PageInfo.PageSize = 100
		resp.PageInfo.Results = len(resp.Results)
		resp.PageInfo.Page = 1

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := &config.JellyseerrConfig{
		Enabled:      true,
		URL:          server.URL,
		APIKey:       config.Secret("test-api-key"),
		ExcludeUsers: []string{"werner", "excluded@example.com"},
	}

	recipients, err := FetchRecipients(context.Background(), cfg, server.Client())
	require.NoError(t, err)
	assert.Len(t, recipients, 2)
	assert.Contains(t, recipients, `"Fabian" <fabian@example.com>`)
	assert.Contains(t, recipients, `"Lorenz" <lorenz@example.com>`)
}

func TestFetchRecipients_Pagination(t *testing.T) {
	pageCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageCount++
		skip := r.URL.Query().Get("skip")

		var users []User
		if skip == "0" {
			users = []User{
				{ID: 1, Email: "user1@example.com", DisplayName: "User 1"},
			}
		} else {
			users = []User{
				{ID: 2, Email: "user2@example.com", DisplayName: "User 2"},
			}
		}

		resp := UsersResponse{Results: users}
		resp.PageInfo.Pages = 2
		resp.PageInfo.PageSize = 1
		resp.PageInfo.Results = 2
		if skip == "0" {
			resp.PageInfo.Page = 1
		} else {
			resp.PageInfo.Page = 2
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := &config.JellyseerrConfig{
		Enabled: true,
		URL:     server.URL,
		APIKey:  config.Secret("test-key"),
	}

	recipients, err := FetchRecipients(context.Background(), cfg, server.Client())
	require.NoError(t, err)
	assert.Equal(t, 2, pageCount)
	assert.Len(t, recipients, 2)
	assert.Contains(t, recipients, `"User 1" <user1@example.com>`)
	assert.Contains(t, recipients, `"User 2" <user2@example.com>`)
}

func TestFetchRecipients_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	cfg := &config.JellyseerrConfig{
		Enabled: true,
		URL:     server.URL,
		APIKey:  config.Secret("bad-key"),
	}

	_, err := FetchRecipients(context.Background(), cfg, server.Client())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status 401")
}
