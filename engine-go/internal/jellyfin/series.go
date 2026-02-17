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

type SeriesItem struct {
	Name           string
	AdditionDate   time.Time
	ProductionYear int32
	Seasons        map[string]SeasonItem
	TMDBId         int
}

type NewlyAddedSeriesItem struct {
	SeriesName     string
	SeriesID       string
	IsSeriesNew    bool
	NewSeasons     map[string]SeasonItem
	TMDBId         int
	ProductionYear int
	AdditionDate   time.Time
}

func parseSeriesItems(app *app.ApplicationContext, jellyfinItems *[]jellyfinAPI.BaseItemDto) map[string]SeriesItem {
	seriesItems := map[string]SeriesItem{}
	for _, item := range *jellyfinItems {
		if *item.Type == jellyfinAPI.BASEITEMKIND_SERIES {
			seriesItems[*item.Id] = SeriesItem{
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

func updateSeriesWithSeasons(
	jellyfinItems *[]jellyfinAPI.BaseItemDto,
	seriesItem map[string]SeriesItem,
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

func updateSeriesWithEpisode(
	jellyfinItems *[]jellyfinAPI.BaseItemDto,
	seriesItem map[string]SeriesItem,
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

func (client *APIClient) fetchAndParseSeries(
	folderName string,
	app *app.ApplicationContext,
) (map[string]SeriesItem, error) {
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

func (client *APIClient) buildNewlyAddedSeriesList(
	seriesItem map[string]SeriesItem,
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

func (client *APIClient) createNewlyAddedSeriesItem(
	seriesID string,
	series SeriesItem,
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
