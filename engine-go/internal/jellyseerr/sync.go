package jellyseerr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
)

const (
	defaultHTTPTimeout = 30 * time.Second
	maxErrorBodyBytes  = 1024
	pageSize           = 100
)

type User struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	Username         string `json:"username"`
	JellyfinUsername string `json:"jellyfinUsername"`
	DisplayName      string `json:"displayName"`
}

type UsersResponse struct {
	PageInfo struct {
		Pages    int `json:"pages"`
		PageSize int `json:"pageSize"`
		Results  int `json:"results"`
		Page     int `json:"page"`
	} `json:"pageInfo"`
	Results []User `json:"results"`
}

// FetchRecipients retrieves users from Jellyseerr/Overseerr and returns formatted email recipients.
func FetchRecipients(ctx context.Context, cfg *config.JellyseerrConfig, client *http.Client) ([]string, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}

	if client == nil {
		client = &http.Client{Timeout: defaultHTTPTimeout}
	}

	excludeSet := buildExcludeSet(cfg.ExcludeUsers)
	seenEmails := make(map[string]struct{})
	var recipients []string
	skip := 0

	for {
		reqURL := fmt.Sprintf("%s/api/v1/user?take=%d&skip=%d", cfg.URL, pageSize, skip)
		page, err := fetchPage(ctx, client, reqURL, cfg.APIKey.SafeString())
		if err != nil {
			return nil, err
		}

		if len(page.Results) == 0 {
			break
		}

		newRecipients := processUsers(page.Results, excludeSet, seenEmails)
		recipients = append(recipients, newRecipients...)

		skip += len(page.Results)
		if page.PageInfo.Pages > 0 {
			if page.PageInfo.Page >= page.PageInfo.Pages {
				break
			}
		} else if len(page.Results) < pageSize {
			break
		}
	}

	return recipients, nil
}

func buildExcludeSet(excludeUsers []string) map[string]struct{} {
	excludeSet := make(map[string]struct{}, len(excludeUsers))
	for _, u := range excludeUsers {
		trimmed := strings.ToLower(strings.TrimSpace(u))
		if trimmed != "" {
			excludeSet[trimmed] = struct{}{}
		}
	}
	return excludeSet
}

func fetchPage(ctx context.Context, client *http.Client, reqURL, apiKey string) (*UsersResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Jellyseerr API: %w", err)
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query Jellyseerr API at %s: %w", reqURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		return nil, fmt.Errorf("jellyseerr API returned unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var page UsersResponse
	err = json.NewDecoder(resp.Body).Decode(&page)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Jellyseerr API response: %w", err)
	}

	return &page, nil
}

func processUsers(users []User, excludeSet, seenEmails map[string]struct{}) []string {
	var recipients []string

	for _, user := range users {
		rawEmail := strings.TrimSpace(user.Email)
		if rawEmail == "" || !strings.Contains(rawEmail, "@") {
			continue
		}

		normalizedEmail := strings.ToLower(rawEmail)
		if _, exists := seenEmails[normalizedEmail]; exists {
			continue
		}

		if isUserExcluded(user, excludeSet) {
			continue
		}

		seenEmails[normalizedEmail] = struct{}{}

		displayName := strings.TrimSpace(user.DisplayName)
		if displayName == "" {
			displayName = strings.TrimSpace(user.JellyfinUsername)
		}
		if displayName == "" {
			displayName = strings.TrimSpace(user.Username)
		}

		addr := mail.Address{
			Name:    displayName,
			Address: rawEmail,
		}
		recipients = append(recipients, addr.String())
	}

	return recipients
}

func isUserExcluded(user User, excludeSet map[string]struct{}) bool {
	candidates := []string{
		strings.ToLower(strings.TrimSpace(user.Email)),
		strings.ToLower(strings.TrimSpace(user.Username)),
		strings.ToLower(strings.TrimSpace(user.JellyfinUsername)),
		strings.ToLower(strings.TrimSpace(user.DisplayName)),
	}

	for _, cand := range candidates {
		if cand != "" {
			if _, found := excludeSet[cand]; found {
				return true
			}
		}
	}

	return false
}
