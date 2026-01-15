package config

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
	TestSMTPConnection bool
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

type yamlConfiguration struct {
	Log *struct {
		Level  string `yaml:"level,omitempty" validate:"omitempty,oneof=DEBUG INFO WARN ERROR"`
		Format string `yaml:"format, omitempty" validate:"omitempty,oneof=json console"`
	} `yaml:"log,omitempty"`
	Scheduler *struct {
		Cron string `yaml:"cron" validate:"cron"`
	} `yaml:"scheduler,omitempty"`
	Jellyfin struct {
		Url                                 string   `yaml:"url" validate:"required,http_url"`
		ApiToken                            string   `yaml:"api_token" validate:"required"`
		WatchedFilmFolders                  []string `yaml:"watched_film_folders" validate:"required"`
		WatchedSeriesFolders                []string `yaml:"watched_tv_folders" validate:"required"`
		ObservedPeriodDays                  int      `yaml:"observed_period_days" validate:"required,numeric"`
		IgnoreItemsAddedAfterLastNewsletter *bool    `yaml:"ignore_item_added_before_last_newsletter,omitempty" validate:"omitempty,boolean"`
	} `yaml:"jellyfin"            validate:"required"`
	Tmdb struct {
		ApiKey string `yaml:"api_key" validate:"required,jwt"`
	} `yaml:"tmdb"                validate:"required"`
	EmailTemplate struct {
		Theme                   string `yaml:"theme,omitempty" validate:"omitempty,oneof=classic"`
		Language                string `yaml:"language" validate:"required,alpha"`
		Subject                 string `yaml:"subject"  validate:"required"`
		Title                   string `yaml:"title"   validate:"required"`
		Subtitle                string `yaml:"subtitle"   validate:"required"`
		JellyfinURL             string `yaml:"jellyfin_url,omitempty" validate:"omitempty,url"`
		UnsubscribeEmail        string `yaml:"unsubscribe_email,omitempty" validate:"omitempty,email"`
		JellyfinOwnerName       string `yaml:"jellyfin_owner_name,omitempty"`
		DisplayOverviewMaxItems *int   `yaml:"display_overview_max_items,omitempty" validate:"omitempty,numeric,min=-1"`
		SortMode                string `yaml:"sort_mode,omitempty" validate:"omitempty,oneof=date_desc date_asc name_asc name_desc"`
	} `yaml:"email_template"      validate:"required"`
	Email struct {
		SmtpServer     string `yaml:"smtp_server" validate:"required,hostname|ip"`
		SmtpPort       int    `yaml:"smtp_port" validate:"required,numeric,min=1,max=65535"`
		SmtpUsername   string `yaml:"smtp_username" validate:"required"`
		SmtpPassword   string `yaml:"smtp_password" validate:"required"`
		SmtpSenderName string `yaml:"smtp_sender_email" validate:"required"`
		SmtpTlsType    string `yaml:"smtp_tls_type,omitempty" validate:"omitempty,oneof=TLS STARTTLS"`
	} `yaml:"email"               validate:"required"`
	DryRun *struct {
		Enabled            bool   `yaml:"enabled" validate:"boolean"`
		TestSMTPConnection *bool  `yaml:"test_smtp_connection,omitempty" validate:"omitempty,boolean"`
		OutputDirectory    string `yaml:"output_directory,omitempty" validate:"omitempty,dirpath"`
		OutputFilename     string `yaml:"output_filename,omitempty" `
		IncludeMetadata    *bool  `yaml:"include_metadata,omitempty" validate:"omitempty,boolean"`
		SaveEmailData      *bool  `yaml:"save_email_data,omitempty" validate:"omitempty,boolean"`
	} `yaml:"dry-run,omitempty"`
	Recipients []string `yaml:"recipients"          validate:"required"`
}
