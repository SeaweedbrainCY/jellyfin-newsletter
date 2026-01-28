package jellyfin

import (
	"context"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/sj14/jellyfin-go/api"
	"go.uber.org/zap"
)

const ITEM_TYPE_SEASON = "season"
const ITEM_TYPE_SERIES = "series"
const ITEM_TYPE_EPISODE = "episode"

type EpisodeItem struct {
	Name           string
	AdditionDate   time.Time
	ProductionYear int32
	EpisodeNumber  int32
}

type SeasonItem struct {
	SeasonNumber   int32
	Name           string
	AdditionDate   time.Time
	ProductionYear int32
	Episodes       map[string]EpisodeItem
	TMDBId         int
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
	NewEpisodes    map[string]EpisodeItem
	TMDBId         int
	ProductionYear int
	AdditionDate   time.Time
}

func (client *APIClient) getAllJellyfinItemsFromFolder(
	folderID string,
	app *app.ApplicationContext,
) (*[]api.BaseItemDto, error) {
	items, httpResponse, err := client.ItemsAPI.GetItems(context.Background()).
		Recursive(true).
		ParentId(folderID).
		Fields([]api.ItemFields{"DateCreated", "ProviderIds", "Id", "Name", "ProductionYear", "IndexNumber", "SeriesId", "Type", "SeasonId"}).
		Execute()
	if err != nil {
		logHTTPResponseError(httpResponse, err, app)
		return nil, err
	}
	return &items.Items, nil
}

func parseSeriesItems(jellyfinItems *[]api.BaseItemDto, app *app.ApplicationContext) map[string]SeriesItem {
	seriesItems := map[string]SeriesItem{}
	for _, item := range *jellyfinItems {
		if *item.Type == api.BASEITEMKIND_SERIES {
			seriesItems[*item.Id] = SeriesItem{
				Name:           OrDefault(item.Name, "Unknown"),
				AdditionDate:   OrDefault(item.DateCreated, time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)),
				ProductionYear: OrDefault(item.ProductionYear, 0),
				Seasons:        map[string]SeasonItem{},
				TMDBId:         getTMDBIDIfExist(&item),
			}
		}
	}
	return seriesItems
}

func updateSeriesWithSeasons(
	jellyfinItems *[]api.BaseItemDto,
	seriesItem map[string]SeriesItem,
	app *app.ApplicationContext,
) {
	for _, item := range *jellyfinItems {
		if *item.Type == api.BASEITEMKIND_SEASON {
			if !item.SeriesId.IsSet() {
				app.Logger.Warn("A season item is ignored because it has no series ID.", zap.String("itemID", *item.Id))
				continue
			}
			seasonItem := SeasonItem{
				Name:           OrDefault(item.Name, "Unknown"),
				AdditionDate:   OrDefault(item.DateCreated, time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)),
				ProductionYear: OrDefault(item.ProductionYear, 0),
				TMDBId:         getTMDBIDIfExist(&item),
				SeasonNumber:   OrDefault(item.IndexNumber, 0),
				Episodes:       map[string]EpisodeItem{},
			}
			if _, ok := seriesItem[*item.SeriesId.Get()]; !ok {
				app.Logger.Warn(
					"A season item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					zap.String("itemID", *item.Id),
					zap.String("seriesID", *item.SeriesId.Get()),
				)
				continue
			}
			seriesItem[*item.SeriesId.Get()].Seasons[*item.Id] = seasonItem
		}
	}
}

func updateSeriesWithEpisode(
	jellyfinItems *[]api.BaseItemDto,
	seriesItem map[string]SeriesItem,
	app *app.ApplicationContext,
) {
	for _, item := range *jellyfinItems {
		if *item.Type == api.BASEITEMKIND_EPISODE && item.LocationType.IsSet() &&
			*item.LocationType.Get() == api.LOCATIONTYPE_FILE_SYSTEM {
			if !item.SeriesId.IsSet() || !item.SeasonId.IsSet() {
				app.Logger.Warn(
					"An episode item is ignored because it has no series ID or season ID.",
					zap.String("itemID", *item.Id),
					zap.String("seasonID", OrDefault(item.SeasonId, "")),
					zap.String("seriesID", OrDefault(item.SeriesId, "")),
				)
				continue
			}
			if _, ok := seriesItem[*item.SeriesId.Get()]; !ok {
				app.Logger.Warn(
					"An episode item is ignored because it belongs to a Series that doesn't exist in Jellyfin's API response.",
					zap.String("itemID", *item.Id),
					zap.String("seriesID", *item.SeriesId.Get()),
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
				Name:           OrDefault(item.Name, ""),
				AdditionDate:   OrDefault(item.DateCreated, time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC)),
				ProductionYear: OrDefault(item.ProductionYear, 0),
				EpisodeNumber:  OrDefault(item.IndexNumber, 0),
			}

			seriesItem[*item.SeriesId.Get()].Seasons[*item.SeasonId.Get()].Episodes[*item.Id] = episodeItem
		}
	}
}

func (client *APIClient) getNewlyAddedSeriesByFolder(
	folderName string,
	app *app.ApplicationContext,
) (*[]NewlyAddedSeriesItem, error) {
	minimumAdditionDate := time.Now().AddDate(0, 0, app.Config.Jellyfin.ObservedPeriodDays*-1)
	app.Logger.Debug(
		"Searching for recently added series.",
		zap.String("FolderName", folderName),
		zap.String("StartAdditionDate", minimumAdditionDate.String()),
	)

	folderID, err := client.GetRootFolderIDByName(folderName, app)
	if err != nil {
		return nil, err
	}

	jellyfinItems, err := client.getAllJellyfinItemsFromFolder(folderID, app)
	if err != nil {
		return nil, err
	}

	seriesItem := parseSeriesItems(jellyfinItems, app)
	updateSeriesWithSeasons(jellyfinItems, seriesItem, app)
	updateSeriesWithEpisode(jellyfinItems, seriesItem, app)

	var newlyAddedSeries = []NewlyAddedSeriesItem{}
	//app.Logger.Info("debug", zap.Any("series", seriesItem))
	for seriesID, series := range seriesItem {
		newSeries := NewlyAddedSeriesItem{
			SeriesName:     series.Name,
			SeriesID:       seriesID,
			NewSeasons:     map[string]SeasonItem{},
			NewEpisodes:    map[string]EpisodeItem{},
			TMDBId:         series.TMDBId,
			ProductionYear: int(series.ProductionYear),
			AdditionDate:   series.AdditionDate,
		}
		if series.AdditionDate.After(minimumAdditionDate) {
			newSeries.IsSeriesNew = true
			newlyAddedSeries = append(newlyAddedSeries, newSeries)
			continue
		}
		newSeries.IsSeriesNew = false
		for seasonID, season := range series.Seasons {
			if season.AdditionDate.After(minimumAdditionDate) {
				newSeries.NewSeasons[seasonID] = season
				continue
			}
			for episodeID, episode := range season.Episodes {
				if episode.AdditionDate.After(minimumAdditionDate) {
					newSeries.NewEpisodes[episodeID] = episode
				}
			}
		}
		if len(newSeries.NewEpisodes) != 0 || len(newSeries.NewSeasons) != 0 {
			newlyAddedSeries = append(newlyAddedSeries, newSeries)
		}
	}
	return &newlyAddedSeries, nil
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
