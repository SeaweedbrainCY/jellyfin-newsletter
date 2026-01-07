package context

import (
	"fmt"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loadContextFromReader(r io.Reader) (*Context, error) {
yamlParsedConfig := &yamlConfiguration{}

	validate := validator.New()

	decoder := yaml.NewDecoder(
		r,
		yaml.Validator(validate),
		yaml.Strict(),
	)

	err := decoder.Decode(yamlParsedConfig)
	if err != nil {
		yaml.FormatError(err, true, true)
		return nil, fmt.Errorf("Failed to decode configuration file: %w", err)
	}

	context := &Context{}
	config := &Configuration{}

	config.Log = LogConfig{
		Level:  "INFO",
		Format: "console",
	}

	if yamlParsedConfig.Log != nil {
		if yamlParsedConfig.Log.Format != "" {
			config.Log.Format = yamlParsedConfig.Log.Format
		}
		if yamlParsedConfig.Log.Level != "" {
			config.Log.Level = yamlParsedConfig.Log.Level
		}
	}

	context.Logger, err = initializeLogger(&config.Log)
	if err != nil {
		return nil, err
	}

	if yamlParsedConfig.Scheduler != nil {
		config.Scheduler.Enabled = true
		config.Scheduler.CronExpr = yamlParsedConfig.Scheduler.Cron
	} else {
		config.Scheduler.Enabled = false
	}

	config.Jellyfin = JellyfinConfig{
		Url:                                 yamlParsedConfig.Jellyfin.Url,
		ApiKey:                              yamlParsedConfig.Jellyfin.ApiToken,
		WatchedFilmFolders:                  yamlParsedConfig.Jellyfin.WatchedFilmFolders,
		WatchedSeriesFolders:                yamlParsedConfig.Jellyfin.WatchedSeriesFolders,
		ObservedPeriodDays:                  yamlParsedConfig.Jellyfin.ObservedPeriodDays,
		IgnoreItemsAddedAfterLastNewsletter: false,
	}
	context.Logger.Info("Jellyfin URL", zap.String("url", config.Jellyfin.Url))
	if yamlParsedConfig.Jellyfin.IgnoreItemsAddedAfterLastNewsletter != nil && *yamlParsedConfig.Jellyfin.IgnoreItemsAddedAfterLastNewsletter {
		config.Jellyfin.IgnoreItemsAddedAfterLastNewsletter = true
	}

	config.Tmdb.ApiKey = yamlParsedConfig.Tmdb.ApiKey

	config.EmailTemplate = EmailTemplateConfig{
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

	config.SMTP = SmtpConfig{
		Host:       yamlParsedConfig.Email.SmtpServer,
		Port:       yamlParsedConfig.Email.SmtpPort,
		Username:   yamlParsedConfig.Email.SmtpUsername,
		Password:   yamlParsedConfig.Email.SmtpPassword,
		SenderName: yamlParsedConfig.Email.SmtpSenderName,
		TlsType:    "STARTTLS",
	}

	if yamlParsedConfig.Email.SmtpTlsType != "" {
		config.SMTP.TlsType = yamlParsedConfig.Email.SmtpTlsType
	}

	config.DryRun.Enabled = false
	if yamlParsedConfig.DryRun != nil && yamlParsedConfig.DryRun.Enabled {
		config.DryRun = DryRunConfig{
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
	context.Config = config
	return context, nil
}

func LoadContext(configPath string) (*Context, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read configuration file: %w", err)
	}
	return loadContextFromReader(file)
	
}

func initializeLogger(logConfiguration *LogConfig) (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logFormat := "console"

	if logConfiguration.Format == "json" || logConfiguration.Format == "console" {
		logFormat = logConfiguration.Format
	}

	var logLevel zap.AtomicLevel
	switch logConfiguration.Level {
	case "DEBUG":
		logLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "INFO":
		logLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "WARN":
		logLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "ERROR":
		logLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		logLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	config := zap.Config{
		Level:            logLevel,
		Development:      false,
		Encoding:         logFormat,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("Error while building logger: %w", err)
	}
	return logger, nil

}
