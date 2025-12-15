package config

import (
	"os"

	"github.com/goccy/go-yaml"
	"go.uber.org/zap"
)

func LoadConfiguration(configPath string, logger *zap.Logger) (*Configuration, error) {

	file, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal("Failed to read configuration file", zap.Error(err))
		return nil, err
	}

	yamlParsedConfig := &yamlConfiguration{}
	err = yaml.Unmarshal(file, yamlParsedConfig)
	if err != nil {
		logger.Fatal("Failed to parse configuration", zap.Error(err))
		return nil, err
	}

	config := &Configuration{}
	config.Logger = logger

	isDebug := false
	if yamlParsedConfig.Debug != nil && *yamlParsedConfig.Debug {
		isDebug = true
	}
	config.Debug = isDebug

	if yamlParsedConfig.Scheduler != nil {
		config.Scheduler.Enabled = true
		config.Scheduler.CronExpr = yamlParsedConfig.Scheduler.Cron
	} else {
		config.Scheduler.Enabled = false
	}

	config.Jellyfin = jellyfinConfig{
		Url:                                 yamlParsedConfig.Jellyfin.Url,
		ApiKey:                              yamlParsedConfig.Jellyfin.ApiKey,
		WatchedFilmFolders:                  yamlParsedConfig.Jellyfin.WatchedFilmFolders,
		WatchedSeriesFolders:                yamlParsedConfig.Jellyfin.WatchedSeriesFolders,
		ObservedPeriodDays:                  yamlParsedConfig.Jellyfin.ObservedPeriodDays,
		IgnoreItemsAddedAfterLastNewsletter: false,
	}
	if yamlParsedConfig.Jellyfin.IgnoreItemsAddedAfterLastNewsletter != nil && *yamlParsedConfig.Jellyfin.IgnoreItemsAddedAfterLastNewsletter {
		config.Jellyfin.IgnoreItemsAddedAfterLastNewsletter = true
	}

	config.Tmdb.ApiKey = yamlParsedConfig.Tmdb.ApiKey

	config.EmailTemplate = emailTemplateConfig{
		Language:                yamlParsedConfig.EmailTemplate.Language,
		Subject:                 yamlParsedConfig.EmailTemplate.Subject,
		Title:                   yamlParsedConfig.EmailTemplate.Title,
		Subtitle:                yamlParsedConfig.EmailTemplate.Subtitle,
		JellyfinURL:             yamlParsedConfig.EmailTemplate.JellyfinURL,
		UnsubscribeEmail:        yamlParsedConfig.EmailTemplate.UnsubscribeEmail,
		JellyfinOwnerName:       yamlParsedConfig.EmailTemplate.JellyfinOwnerName,
		Theme:                   "classic",
		DisplayOverviewMaxItems: 10,
		SortMode:                "date_desc",
		AvailableLanguages:      []string{"en", "fr", "he", "ca", "es", "it"},
	}
	if yamlParsedConfig.EmailTemplate.Theme != "" {
		config.EmailTemplate.Theme = yamlParsedConfig.EmailTemplate.Theme
	}
	if yamlParsedConfig.EmailTemplate.DisplayOverviewMaxItems != nil {
		config.EmailTemplate.DisplayOverviewMaxItems = *yamlParsedConfig.EmailTemplate.DisplayOverviewMaxItems
	}
	if yamlParsedConfig.EmailTemplate.SortMode != "" {
		config.EmailTemplate.SortMode = yamlParsedConfig.EmailTemplate.SortMode
	}

	config.SMTP = smtpConfig{
		Host:       yamlParsedConfig.SMTP.Host,
		Port:       yamlParsedConfig.SMTP.Port,
		Username:   yamlParsedConfig.SMTP.Username,
		Password:   yamlParsedConfig.SMTP.Password,
		SenderName: yamlParsedConfig.SMTP.SenderName,
		TlsType:    "STARTTLS",
	}

	if yamlParsedConfig.SMTP.TlsType != "" {
		config.SMTP.TlsType = yamlParsedConfig.SMTP.TlsType
	}

	config.DryRun.Enabled = false
	if yamlParsedConfig.DryRun != nil && yamlParsedConfig.DryRun.Enabled {
		config.DryRun = dryRunConfig{
			Enabled:            true,
			TestSTMPConnection: false,
			OutputDirectory:    "./previews/",
			OutputFilename:     "newsletter_{date}_{time}.html",
			IncludeMetadata:    true,
			SaveEmailData:      true,
		}
		if yamlParsedConfig.DryRun.TestSTMPConnection != nil && *yamlParsedConfig.DryRun.TestSTMPConnection {
			config.DryRun.TestSTMPConnection = true
		}
		if yamlParsedConfig.DryRun.OutputDirectory != "" {
			config.DryRun.OutputDirectory = yamlParsedConfig.DryRun.OutputDirectory
		}
		if yamlParsedConfig.DryRun.OutputFilename != "" {
			config.DryRun.OutputFilename = yamlParsedConfig.DryRun.OutputFilename
		}
		if yamlParsedConfig.DryRun.IncludeMetadata != nil && *yamlParsedConfig.DryRun.IncludeMetadata {
			config.DryRun.IncludeMetadata = true
		}
		if yamlParsedConfig.DryRun.SaveEmailData != nil && *yamlParsedConfig.DryRun.SaveEmailData {
			config.DryRun.SaveEmailData = true
		}
	}
	config.EmailRecipients = yamlParsedConfig.Recipients

	return config, nil
}
