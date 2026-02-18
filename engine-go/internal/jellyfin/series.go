package jellyfin

import (
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	jellyfinAPI "github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

type EpisodeItem struct {
	Name          string
	AdditionDate  time.Time
	EpisodeNumber int32
}

type SeasonItem struct {
	SeasonNumber int32
	Name         string
	AdditionDate time.Time
	Episodes     map[string]EpisodeItem
	IsSeasonNew  bool
}

type seriesItem struct {
	Name           string
	AdditionDate   time.Time
	ProductionYear int32
	Seasons        map[string]SeasonItem
	TMDBId         string
}

type NewlyAddedSeriesItem struct {
	SeriesName     string
	SeriesID       string
	IsSeriesNew    bool
	NewSeasons     map[string]SeasonItem
	TMDBId         string
	ProductionYear int
	AdditionDate   time.Time
}

// parseSeriesItems scans a slice of Jellyfin BaseItemDto and extracts
// all items of type SERIES into a map keyed by the series ID.
// For each series it builds a `seriesItem` with name, addition date,
// production year, an empty seasons map and the TMDB id when present.
// Missing addition dates are logged as warnings because they affect
// newly-added detection.
func parseSeriesItems(app *app.ApplicationContext, jellyfinItems *[]jellyfinAPI.BaseItemDto) map[string]seriesItem {
	seriesItems := map[string]seriesItem{}
	for _, item := range *jellyfinItems {
		if *item.Type == jellyfinAPI.BASEITEMKIND_SERIES {
			seriesItems[*item.Id] = seriesItem{
				Name:           OrDefault(item.Name, "Unknown"),
				AdditionDate:   OrDefault(item.DateCreated, time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)),
				ProductionYear: OrDefault(item.ProductionYear, 0),
				Seasons:        map[string]SeasonItem{},
				TMDBId:         getTMDBIDIfExist(&item),
			}
			if seriesItems[*item.Id].AdditionDate.Equal(time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)) {
				app.Logger.Warn(
					"Found a series with no addition date. This can lead to inaccuracy when detecting newly added media.",
					zap.String("Series ID", *item.Id),
					zap.String("Series Name", seriesItems[*item.Id].Name),
				)
			}
		}
	}
	return seriesItems
}

// updateSeriesWithSeasons iterates over Jellyfin items and attaches
// Season entries to their parent `seriesItem` in the provided map.
// Seasons that reference a missing series are ignored and logged.
// If a season lacks an addition date the function emits a warning.
func updateSeriesWithSeasons(
	jellyfinItems *[]jellyfinAPI.BaseItemDto,
	seriesItem map[string]seriesItem,
	app *app.ApplicationContext,
) {
	for _, item := range *jellyfinItems {
		if *item.Type == jellyfinAPI.BASEITEMKIND_SEASON {
			if !item.SeriesId.IsSet() || item.SeriesId.Get() == nil {
				app.Logger.Warn(
					"A season item is ignored because it has no series ID.",
					zap.String("Season ID", *item.Id),
				)
				continue
			}
			seasonItem := SeasonItem{
				Name:         OrDefault(item.Name, "Unknown"),
				AdditionDate: OrDefault(item.DateCreated, time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)),
				SeasonNumber: OrDefault(item.IndexNumber, 0),
				Episodes:     map[string]EpisodeItem{},
			}
			if _, ok := seriesItem[*item.SeriesId.Get()]; !ok {
				app.Logger.Warn(
					"A season item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					zap.String("Season ID", *item.Id),
					zap.String("Season Name", seasonItem.Name),
					zap.String("Not found Series Name", OrDefault(item.SeriesName, "Unknown")),
					zap.String("Not found Series ID", OrDefault(item.SeriesId, "Unknown")),
				)
				continue
			}
			if seasonItem.AdditionDate.Equal(time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)) {
				app.Logger.Warn(
					"Found a season with no addition date. This can lead to inaccuracy when detecting newly added media.",
					zap.String("Season ID", *item.Id),
					zap.String("Season Name", seasonItem.Name),
					zap.String("Series Name", OrDefault(item.SeriesName, "Unknown")),
					zap.String("Series ID", *item.SeriesId.Get()),
				)
			}
			seriesItem[*item.SeriesId.Get()].Seasons[*item.Id] = seasonItem
		}
	}
}

// updateSeriesWithEpisode finds episode items (file-system location)
// and attaches them into the corresponding `SeasonItem.Episodes` map
// under the correct `seriesItem`. Episodes referencing missing series
// or seasons are ignored and logged. Episodes without addition dates
// also produce warnings because they affect new-content detection.
func updateSeriesWithEpisode(
	jellyfinItems *[]jellyfinAPI.BaseItemDto,
	seriesItem map[string]seriesItem,
	app *app.ApplicationContext,
) {
	for _, item := range *jellyfinItems {
		if *item.Type == jellyfinAPI.BASEITEMKIND_EPISODE && item.LocationType.IsSet() &&
			item.LocationType.Get() != nil &&
			*item.LocationType.Get() == jellyfinAPI.LOCATIONTYPE_FILE_SYSTEM {
			if !item.SeriesId.IsSet() || !item.SeasonId.IsSet() || item.SeriesId.Get() == nil ||
				item.SeasonId.Get() == nil {
				app.Logger.Warn(
					"An episode item is ignored because it has no series ID or season ID.",
					zap.String("Episode ID", *item.Id),
					zap.String("Episode Name", OrDefault(item.Name, "Unknown")),
					zap.String("Expected Series Name", OrDefault(item.SeriesName, "Unknown")),
					zap.String("Expected Series ID", OrDefault(item.SeriesId, "Unknown")),
					zap.String("Expected Season Name", OrDefault(item.SeasonName, "Unknown")),
					zap.String("Expected Season ID", OrDefault(item.SeasonId, "Unknown")),
				)
				continue
			}
			if _, ok := seriesItem[*item.SeriesId.Get()]; !ok {
				app.Logger.Warn(
					"An episode item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					zap.String("Episode ID", *item.Id),
					zap.String("Episode Name", OrDefault(item.Name, "Unknown")),
					zap.String("Expected Series Name", OrDefault(item.SeriesName, "Unknown")),
					zap.String("Expected Series ID", OrDefault(item.SeriesId, "Unknown")),
					zap.String("Expected Season Name", OrDefault(item.SeasonName, "Unknown")),
					zap.String("Expected Season ID", OrDefault(item.SeasonId, "Unknown")),
				)
				continue
			}
			if _, ok := seriesItem[*item.SeriesId.Get()].Seasons[*item.SeasonId.Get()]; !ok {
				app.Logger.Warn(
					"An episode item is ignored because it belongs to a Seasons that doesn't exist in Jellyfin's API response.",
					zap.String("itemID", *item.Id),
					zap.String("seasonID", *item.SeasonId.Get()),
				)
				continue
			}
			episodeItem := EpisodeItem{
				Name:          OrDefault(item.Name, "Unknown"),
				AdditionDate:  OrDefault(item.DateCreated, time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)),
				EpisodeNumber: OrDefault(item.IndexNumber, 0),
			}
			if episodeItem.AdditionDate.Equal(time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)) {
				app.Logger.Warn(
					"Found an episode with no addition date. This can lead to inaccuracy when detecting newly added media.",
					zap.String("Episode ID", *item.Id),
					zap.String("Episode Name", episodeItem.Name),
					zap.String("Season Name", OrDefault(item.SeasonName, "Unknown")),
					zap.String("Season ID", *item.SeasonId.Get()),
					zap.String("Series Name", OrDefault(item.SeriesName, "Unknown")),
					zap.String("Series ID", *item.SeriesId.Get()),
				)
			}

			seriesItem[*item.SeriesId.Get()].Seasons[*item.SeasonId.Get()].Episodes[*item.Id] = episodeItem
		}
	}
}

// getNewlyAddedSeriesByFolder fetches series metadata for the given
// Jellyfin folder, computes the cutoff date based on the configured
// observed period, and returns a slice of `NewlyAddedSeriesItem`
// representing series, new seasons or episodes added after the cutoff.
func (client *APIClient) getNewlyAddedSeriesByFolder(
	folderName string,
	app *app.ApplicationContext,
) (*[]NewlyAddedSeriesItem, error) {
	minimumAdditionDate := time.Now().AddDate(0, 0, app.Config.Jellyfin.ObservedPeriodDays*-1-1)
	app.Logger.Debug(
		"Searching for recently added series.",
		zap.String("FolderName", folderName),
		zap.String("StartAdditionDate", minimumAdditionDate.String()),
	)

	seriesItem, err := client.fetchAndParseSeries(folderName, app)
	if err != nil {
		return nil, err
	}

	newlyAddedSeries := client.buildNewlyAddedSeriesList(seriesItem, minimumAdditionDate)
	return &newlyAddedSeries, nil
}

// fetchAndParseSeries resolves the folder ID by name, retrieves all
// items for that folder from the Items API, and builds a structured
// map of `seriesItem` populated with seasons and episodes.
// Returns the resulting map or an error encountered while calling
// the Items API.
func (client *APIClient) fetchAndParseSeries(
	folderName string,
	app *app.ApplicationContext,
) (map[string]seriesItem, error) {
	folderID, err := client.ItemsAPI.GetRootFolderIDByName(folderName, app)
	if err != nil {
		return nil, err
	}

	jellyfinItems, err := client.ItemsAPI.GetAllItemsByFolderID(folderID, app)
	if err != nil {
		return nil, err
	}

	seriesItem := parseSeriesItems(app, jellyfinItems)
	updateSeriesWithSeasons(jellyfinItems, seriesItem, app)
	updateSeriesWithEpisode(jellyfinItems, seriesItem, app)

	return seriesItem, nil
}

// buildNewlyAddedSeriesList walks the parsed series map and converts
// each `seriesItem` into a `NewlyAddedSeriesItem` using
// `createNewlyAddedSeriesItem`. Only series that are entirely new or
// contain new seasons/episodes (after `minimumAdditionDate`) are
// included in the returned slice.
func (client *APIClient) buildNewlyAddedSeriesList(
	seriesItem map[string]seriesItem,
	minimumAdditionDate time.Time,
) []NewlyAddedSeriesItem {
	var newlyAddedSeries []NewlyAddedSeriesItem

	for seriesID, series := range seriesItem {
		newSeries := client.createNewlyAddedSeriesItem(seriesID, series, minimumAdditionDate)

		if newSeries.IsSeriesNew || newSeries.NewSeasons != nil {
			newlyAddedSeries = append(newlyAddedSeries, newSeries)
		}
	}

	return newlyAddedSeries
}

// createNewlyAddedSeriesItem builds a `NewlyAddedSeriesItem` from a
// `seriesItem`. If the series addition date is after the cutoff the
// series is marked as new. Otherwise the function scans seasons to
// detect newly added seasons or episodes and populates `NewSeasons`.
func (client *APIClient) createNewlyAddedSeriesItem(
	seriesID string,
	series seriesItem,
	minimumAdditionDate time.Time,
) NewlyAddedSeriesItem {
	newSeries := NewlyAddedSeriesItem{
		SeriesName:     series.Name,
		SeriesID:       seriesID,
		TMDBId:         series.TMDBId,
		ProductionYear: int(series.ProductionYear),
		AdditionDate:   series.AdditionDate,
	}

	if series.AdditionDate.After(minimumAdditionDate) {
		newSeries.IsSeriesNew = true
		return newSeries
	}

	newSeries.IsSeriesNew = false
	newSeries.NewSeasons = client.findNewSeasons(series.Seasons, minimumAdditionDate)

	return newSeries
}

// findNewSeasons iterates over seasons and returns a map of seasons
// that are newly added or contain newly added episodes relative to
// `minimumAdditionDate`. The returned map is nil when no new seasons
// are found.
func (client *APIClient) findNewSeasons(
	seasons map[string]SeasonItem,
	minimumAdditionDate time.Time,
) map[string]SeasonItem {
	var newSeasons map[string]SeasonItem

	for seasonID, season := range seasons {
		newSeason := client.processSeasonForNewContent(season, minimumAdditionDate)

		if newSeason.IsSeasonNew || newSeason.Episodes != nil {
			if newSeasons == nil {
				newSeasons = map[string]SeasonItem{}
			}
			newSeasons[seasonID] = newSeason
		}
	}

	return newSeasons
}

// processSeasonForNewContent returns a `SeasonItem` describing whether
// the season itself is new (based on addition date) or contains newly
// added episodes. If the season is new it is marked accordingly and
// returned without episode details; otherwise the episode map is
// scanned and returned when new episodes are present.
func (client *APIClient) processSeasonForNewContent(
	season SeasonItem,
	minimumAdditionDate time.Time,
) SeasonItem {
	newSeason := SeasonItem{
		SeasonNumber: season.SeasonNumber,
		Name:         season.Name,
		AdditionDate: season.AdditionDate,
	}

	if season.AdditionDate.After(minimumAdditionDate) {
		newSeason.IsSeasonNew = true
		return newSeason
	}

	newSeason.IsSeasonNew = false
	newSeason.Episodes = client.findNewEpisodes(season.Episodes, minimumAdditionDate)

	return newSeason
}

// findNewEpisodes filters episodes and returns a map containing only
// those episodes whose addition date is after `minimumAdditionDate`.
// Returns nil when no new episodes are found.
func (client *APIClient) findNewEpisodes(
	episodes map[string]EpisodeItem,
	minimumAdditionDate time.Time,
) map[string]EpisodeItem {
	var newEpisodes map[string]EpisodeItem

	for episodeID, episode := range episodes {
		if episode.AdditionDate.After(minimumAdditionDate) {
			if newEpisodes == nil {
				newEpisodes = map[string]EpisodeItem{}
			}
			newEpisodes[episodeID] = episode
		}
	}

	return newEpisodes
}

// GetNewlyAddedSeries collects newly added series information for all
// configured `WatchedSeriesFolders`. For each configured folder it
// calls `getNewlyAddedSeriesByFolder` and aggregates the results.
func (client *APIClient) GetNewlyAddedSeries(
	app *app.ApplicationContext,
) *[]NewlyAddedSeriesItem {
	var seriesItems = []NewlyAddedSeriesItem{}
	for _, folderName := range app.Config.Jellyfin.WatchedSeriesFolders {
		if items, err := client.getNewlyAddedSeriesByFolder(folderName, app); err == nil {
			seriesItems = append(seriesItems, *items...)
		}
	}
	return &seriesItems
}
