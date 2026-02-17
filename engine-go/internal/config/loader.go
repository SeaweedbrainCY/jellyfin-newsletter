package config

import (
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

func LoadConfig(configPath string) (*Configuration, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}
	return loadConfigFromReader(file)
}

func loadConfigFromReader(r io.Reader) (*Configuration, error) {
	yamlParsedConfig := &yamlConfiguration{}

	if err := parseYaml(r, yamlParsedConfig); err != nil {
		return nil, err
	}

	config := &Configuration{}

	config.Log = buildLogConfig(yamlParsedConfig)
	config.Scheduler = buildSchedulerConfig(yamlParsedConfig)
	config.Jellyfin = buildJellyfinConfig(yamlParsedConfig)
	config.TMDB = buildTMDBConfig(yamlParsedConfig)
	config.EmailTemplate = buildEmailTemplateConfig(yamlParsedConfig)
	config.SMTP = buildSMTPConfig(yamlParsedConfig)
	config.DryRun = buildDryRunConfig(yamlParsedConfig)
	config.EmailRecipients = buildRecipientsConfig(yamlParsedConfig)

	return config, nil
}

func parseYaml(r io.Reader, yamlParsedConfig *yamlConfiguration) error {
	validate := validator.New()

	decoder := yaml.NewDecoder(
		r,
		yaml.Validator(validate),
		yaml.Strict(),
	)

	err := decoder.Decode(yamlParsedConfig)
	if err != nil {
		yaml.FormatError(err, true, true)
		return fmt.Errorf("failed to decode configuration file: %w", err)
	}
	return nil
}

func buildLogConfig(yamlParsedConfig *yamlConfiguration) LogConfig {
	logConfig := LogConfig{
		Level:  "INFO",
		Format: "json",
	}

	if yamlParsedConfig.Log != nil {
		if yamlParsedConfig.Log.Format != "" {
			logConfig.Format = yamlParsedConfig.Log.Format
		}
		if yamlParsedConfig.Log.Level != "" {
			logConfig.Level = yamlParsedConfig.Log.Level
		}
	}
	return logConfig
}

func buildSchedulerConfig(yamlParsedConfig *yamlConfiguration) SchedulerConfig {
	schedulerConfig := SchedulerConfig{
		Enabled: false,
	}
	if yamlParsedConfig.Scheduler != nil {
		schedulerConfig.Enabled = true
		schedulerConfig.CronExpr = yamlParsedConfig.Scheduler.Cron
	}
	return schedulerConfig
}

func buildJellyfinConfig(yamlParsedConfig *yamlConfiguration) JellyfinConfig {
	jellyfinConfig := JellyfinConfig{
		URL:                                 yamlParsedConfig.Jellyfin.URL,
		APIKey:                              yamlParsedConfig.Jellyfin.APIToken,
		WatchedFilmFolders:                  yamlParsedConfig.Jellyfin.WatchedFilmFolders,
		WatchedSeriesFolders:                yamlParsedConfig.Jellyfin.WatchedSeriesFolders,
		ObservedPeriodDays:                  yamlParsedConfig.Jellyfin.ObservedPeriodDays,
		IgnoreItemsAddedAfterLastNewsletter: false,
	}
	if yamlParsedConfig.Jellyfin.IgnoreItemsAddedAfterLastNewsletter != nil &&
		*yamlParsedConfig.Jellyfin.IgnoreItemsAddedAfterLastNewsletter {
		jellyfinConfig.IgnoreItemsAddedAfterLastNewsletter = true
	}
	return jellyfinConfig
}

func buildTMDBConfig(yamlParsedConfig *yamlConfiguration) TMDBConfig {
	return TMDBConfig{
		APIKey: yamlParsedConfig.TMDB.APIKey,
	}
}

func buildEmailTemplateConfig(yamlParsedConfig *yamlConfiguration) EmailTemplateConfig {
	const defaultDisplayOverviewMaxItem int = 10

	emailTemplateConfig := EmailTemplateConfig{
		Language:                yamlParsedConfig.EmailTemplate.Language,
		Subject:                 yamlParsedConfig.EmailTemplate.Subject,
		Title:                   yamlParsedConfig.EmailTemplate.Title,
		Subtitle:                yamlParsedConfig.EmailTemplate.Subtitle,
		JellyfinURL:             yamlParsedConfig.EmailTemplate.JellyfinURL,
		UnsubscribeEmail:        yamlParsedConfig.EmailTemplate.UnsubscribeEmail,
		JellyfinOwnerName:       yamlParsedConfig.EmailTemplate.JellyfinOwnerName,
		Theme:                   "classic",
		DisplayOverviewMaxItems: defaultDisplayOverviewMaxItem,
		SortMode:                "date_desc",
		AvailableLanguages:      []string{"en", "fr", "he", "ca", "es", "it"},
	}

	if yamlParsedConfig.EmailTemplate.Theme != "" {
		emailTemplateConfig.Theme = yamlParsedConfig.EmailTemplate.Theme
	}

	if yamlParsedConfig.EmailTemplate.DisplayOverviewMaxItems != nil {
		emailTemplateConfig.DisplayOverviewMaxItems = *yamlParsedConfig.EmailTemplate.DisplayOverviewMaxItems
	}

	if yamlParsedConfig.EmailTemplate.SortMode != "" {
		emailTemplateConfig.SortMode = yamlParsedConfig.EmailTemplate.SortMode
	}
	return emailTemplateConfig
}

func buildSMTPConfig(yamlParsedConfig *yamlConfiguration) SMTPConfig {
	smtpConfig := SMTPConfig{
		Host:       yamlParsedConfig.Email.SMTPServer,
		Port:       yamlParsedConfig.Email.SMTPPort,
		Username:   yamlParsedConfig.Email.SMTPUsername,
		Password:   yamlParsedConfig.Email.SMTPPassword,
		SenderName: yamlParsedConfig.Email.SMTPSenderName,
		TLSType:    "STARTTLS",
	}

	if yamlParsedConfig.Email.SMTPTlsType != "" {
		smtpConfig.TLSType = yamlParsedConfig.Email.SMTPTlsType
	}

	return smtpConfig
}

func buildDryRunConfig(yamlParsedConfig *yamlConfiguration) DryRunConfig {
	dryRunConfig := DryRunConfig{
		Enabled: false,
	}
	if yamlParsedConfig.DryRun != nil && yamlParsedConfig.DryRun.Enabled {
		dryRunConfig = DryRunConfig{
			Enabled: true,
			TestSMTPConnection: yamlParsedConfig.DryRun.TestSMTPConnection != nil &&
				*yamlParsedConfig.DryRun.TestSMTPConnection,
			OutputDirectory: "./previews/",
			OutputFilename:  "newsletter_{date}_{time}.html",
			IncludeMetadata: yamlParsedConfig.DryRun.IncludeMetadata != nil && *yamlParsedConfig.DryRun.IncludeMetadata,
			SaveEmailData:   yamlParsedConfig.DryRun.SaveEmailData != nil && *yamlParsedConfig.DryRun.SaveEmailData,
		}

		if yamlParsedConfig.DryRun.OutputDirectory != "" {
			dryRunConfig.OutputDirectory = yamlParsedConfig.DryRun.OutputDirectory
		}

		if yamlParsedConfig.DryRun.OutputFilename != "" {
			dryRunConfig.OutputFilename = yamlParsedConfig.DryRun.OutputFilename
		}
	}
	return dryRunConfig
}

func buildRecipientsConfig(yamlParsedConfig *yamlConfiguration) []string {
	return yamlParsedConfig.Recipients
}
