package jellyfin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPathExcluded(t *testing.T) {
	tests := []struct {
		name          string
		itemPath      string
		excludedPaths []string
		expected      bool
	}{
		{
			name:          "empty path",
			itemPath:      "",
			excludedPaths: []string{"youtube"},
			expected:      false,
		},
		{
			name:          "empty excluded list",
			itemPath:      "/media/youtube/video.mp4",
			excludedPaths: nil,
			expected:      false,
		},
		{
			name:          "matching path substring",
			itemPath:      "/opt/docker/data/media/youtube/channel/vid.mkv",
			excludedPaths: []string{"youtube"},
			expected:      true,
		},
		{
			name:          "matching path case-insensitive",
			itemPath:      "/media/YouTube-Music/song.mp3",
			excludedPaths: []string{"youtube-music"},
			expected:      true,
		},
		{
			name:          "non-matching path",
			itemPath:      "/media/movies/The Matrix (1999)/matrix.mkv",
			excludedPaths: []string{"youtube", "youtube-music"},
			expected:      false,
		},
		{
			name:          "empty string pattern in list ignored",
			itemPath:      "/media/movies/matrix.mkv",
			excludedPaths: []string{"", "  "},
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPathExcluded(tt.itemPath, tt.excludedPaths)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasMatchingProviderID(t *testing.T) {
	val := "dQw4w9WgXcQ"
	tests := []struct {
		name        string
		providerIDs map[string]*string
		patterns    []string
		expected    bool
	}{
		{
			name:        "nil provider ids",
			providerIDs: nil,
			patterns:    []string{"youtube"},
			expected:    false,
		},
		{
			name:        "empty patterns",
			providerIDs: map[string]*string{"Youtube": &val},
			patterns:    nil,
			expected:    false,
		},
		{
			name:        "matching key",
			providerIDs: map[string]*string{"Youtube": &val},
			patterns:    []string{"youtube"},
			expected:    true,
		},
		{
			name:        "matching value",
			providerIDs: map[string]*string{"Custom": &val},
			patterns:    []string{"dQw4w9WgXcQ"},
			expected:    true,
		},
		{
			name:        "no match",
			providerIDs: map[string]*string{"Tmdb": new("12345")},
			patterns:    []string{"youtube"},
			expected:    false,
		},
		{
			name:        "empty pattern ignored",
			providerIDs: map[string]*string{"Tmdb": new("12345")},
			patterns:    []string{"", "  "},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasMatchingProviderID(tt.providerIDs, tt.patterns)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsGenreOrTagExcluded(t *testing.T) {
	tests := []struct {
		name           string
		genres         []string
		tags           []string
		excludedGenres []string
		expected       bool
	}{
		{
			name:           "empty excluded genres",
			genres:         []string{"Anime", "Action"},
			tags:           []string{"Shounen"},
			excludedGenres: nil,
			expected:       false,
		},
		{
			name:           "matching genre case-insensitive",
			genres:         []string{"Animation", "Anime"},
			tags:           nil,
			excludedGenres: []string{"anime"},
			expected:       true,
		},
		{
			name:           "matching tag",
			genres:         []string{"Animation"},
			tags:           []string{"Japanese Anime"},
			excludedGenres: []string{"anime"},
			expected:       true,
		},
		{
			name:           "no match",
			genres:         []string{"Sci-Fi", "Action"},
			tags:           []string{"Cyberpunk"},
			excludedGenres: []string{"Anime", "Music"},
			expected:       false,
		},
		{
			name:           "empty patterns ignored",
			genres:         []string{"Sci-Fi"},
			tags:           nil,
			excludedGenres: []string{"", "  "},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsGenreOrTagExcluded(tt.genres, tt.tags, tt.excludedGenres)
			assert.Equal(t, tt.expected, result)
		})
	}
}
