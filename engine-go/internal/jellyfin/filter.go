package jellyfin

import "strings"

// IsPathExcluded reports whether itemPath matches any excluded path pattern (case-insensitive substring).
func IsPathExcluded(itemPath string, excludedPaths []string) bool {
	if itemPath == "" || len(excludedPaths) == 0 {
		return false
	}
	lower := strings.ToLower(itemPath)
	for _, pattern := range excludedPaths {
		p := strings.ToLower(strings.TrimSpace(pattern))
		if p != "" && strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// HasMatchingProviderID reports whether any key or value in providerIDs matches any pattern (case-insensitive substring).
func HasMatchingProviderID(providerIDs map[string]*string, patterns []string) bool {
	if len(providerIDs) == 0 || len(patterns) == 0 {
		return false
	}
	for _, pattern := range patterns {
		p := strings.ToLower(strings.TrimSpace(pattern))
		if p == "" {
			continue
		}
		for k, v := range providerIDs {
			if strings.Contains(strings.ToLower(k), p) {
				return true
			}
			if v != nil && strings.Contains(strings.ToLower(*v), p) {
				return true
			}
		}
	}
	return false
}

// IsGenreOrTagExcluded reports whether any genre or tag matches any excluded pattern (case-insensitive substring).
func IsGenreOrTagExcluded(genres []string, tags []string, excludedGenres []string) bool {
	if len(excludedGenres) == 0 {
		return false
	}
	for _, excluded := range excludedGenres {
		e := strings.ToLower(strings.TrimSpace(excluded))
		if e == "" {
			continue
		}
		if containsMatchingString(genres, e) || containsMatchingString(tags, e) {
			return true
		}
	}
	return false
}

func containsMatchingString(list []string, pattern string) bool {
	for _, item := range list {
		if strings.Contains(strings.ToLower(strings.TrimSpace(item)), pattern) {
			return true
		}
	}
	return false
}
