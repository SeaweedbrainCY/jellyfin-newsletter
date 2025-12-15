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
		Level  string `yaml:"level,omitempty"`
		Format string `yaml:"format, omitempty"`
	} `yaml:"log,omitempty"`
	Scheduler *struct {
		Cron string `yaml:"cron"`
	} `yaml:"scheduler,omitempty"`
	Jellyfin struct {
		Url                                 string   `yaml:"url"`
		ApiKey                              string   `yaml:"api_token"`
		WatchedFilmFolders                  []string `yaml:"watched_film_folders"`
		WatchedSeriesFolders                []string `yaml:"watched_tv_folders"`
		ObservedPeriodDays                  int      `yaml:"observed_period_days"`
		IgnoreItemsAddedAfterLastNewsletter *bool    `yaml:"ignore_item_added_before_last_newsletter,omitempty"`
	} `yaml:"jellyfin"`
	Tmdb struct {
		ApiKey string `yaml:"api_key"`
	} `yaml:"tmdb"`
	EmailTemplate struct {
		Theme                   string `yaml:"theme,omitempty"`
		Language                string `yaml:"language,"`
		Subject                 string `yaml:"subject"`
		Title                   string `yaml:"title"`
		Subtitle                string `yaml:"subtitle"`
		JellyfinURL             string `yaml:"jellyfin_url,omitempty"`
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
