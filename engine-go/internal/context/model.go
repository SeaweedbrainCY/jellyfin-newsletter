package context

import (
	"go.uber.org/zap"
)

type LogConfig struct {
	Level  string
	Format string
}

type SchedulerConfig struct {
	Enabled  bool
	CronExpr string
}

type JellyfinConfig struct {
	Url                                 string
	ApiKey                              string
	WatchedFilmFolders                  []string
	WatchedSeriesFolders                []string
	ObservedPeriodDays                  int
	IgnoreItemsAddedAfterLastNewsletter bool
}

type TmdbConfig struct {
	ApiKey string
}

type EmailTemplateConfig struct {
	Theme                   string
	Language                string
	Subject                 string
	Title                   string
	Subtitle                string
	JellyfinURL             string
	UnsubscribeEmail        string
	JellyfinOwnerName       string
	DisplayOverviewMaxItems int
	SortMode                string
	AvailableLanguages      []string
}

type SmtpConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	SenderName string
	TlsType    string
}

type DryRunConfig struct {
	Enabled            bool
	TestSTMPConnection bool
	OutputDirectory    string
	OutputFilename     string
	IncludeMetadata    bool
	SaveEmailData      bool
}

type Configuration struct {
	Log             LogConfig
	EmailRecipients []string
	Scheduler       SchedulerConfig
	Jellyfin        JellyfinConfig
	Tmdb            TmdbConfig
	EmailTemplate   EmailTemplateConfig
	SMTP            SmtpConfig
	DryRun          DryRunConfig
}

type Context struct {
	Config *Configuration
	Logger *zap.Logger
}

type yamlConfiguration struct {
	Log *struct {
		Level  string `yaml:"level,omitempty" validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
		Format string `yaml:"format, omitempty" validate:"omitempty,oneof=json console"`
	} `yaml:"log,omitempty"`
	Scheduler *struct {
		Cron string `yaml:"cron"`
	} `yaml:"scheduler,omitempty"`
	Jellyfin struct {
		Url                                 string   `yaml:"url" validate:"required,http_url"`
		ApiKey                              string   `yaml:"api_token" validate:"required"`
		WatchedFilmFolders                  []string `yaml:"watched_film_folders" validate:"required"`
		WatchedSeriesFolders                []string `yaml:"watched_tv_folders" validate:"required"`
		ObservedPeriodDays                  int      `yaml:"observed_period_days" validate:"required,numeric"`
		IgnoreItemsAddedAfterLastNewsletter *bool    `yaml:"ignore_item_added_before_last_newsletter,omitempty" validate:"boolean"`
	} `yaml:"jellyfin" validate:"required"`
	Tmdb struct {
		ApiKey string `yaml:"api_key" validate:"required,jwt"`
	} `yaml:"tmdb" validate:"required"`
	EmailTemplate struct {
		Theme                   string `yaml:"theme,omitempty"`
		Language                string `yaml:"language" validate:"required,alpha"`
		Subject                 string `yaml:"subject"  validate:"required"`
		Title                   string `yaml:"title"   validate:"required"`
		Subtitle                string `yaml:"subtitle"   validate:"required"`
		JellyfinURL             string `yaml:"jellyfin_url,omitempty" validate:"url"`
		UnsubscribeEmail        string `yaml:"unsubscribe_email,omitempty"`
		JellyfinOwnerName       string `yaml:"jellyfin_owner_name,omitempty"`
		DisplayOverviewMaxItems *int   `yaml:"display_overview_max_items,omitempty"`
		SortMode                string `yaml:"sort_mode,omitempty"`
	} `yaml:"email_template"`
	SMTP struct {
		Host       string `yaml:"smtp_server"`
		Port       int    `yaml:"smtp_port"`
		Username   string `yaml:"smtp_username"`
		Password   string `yaml:"smtp_password"`
		SenderName string `yaml:"smtp_sender_email"`
		TlsType    string `yaml:"smtp_tls_type,omitempty"`
	} `yaml:"email"`
	DryRun *struct {
		Enabled            bool   `yaml:"enabled"`
		TestSTMPConnection *bool  `yaml:"test_smtp_connection,omitempty"`
		OutputDirectory    string `yaml:"output_directory,omitempty"`
		OutputFilename     string `yaml:"output_filename,omitempty"`
		IncludeMetadata    *bool  `yaml:"include_metadata,omitempty"`
		SaveEmailData      *bool  `yaml:"save_email_data,omitempty"`
	} `yaml:"dry_run,omitempty"`
	Recipients []string `yaml:"recipients"`
}
