package config

type schedulerConfig struct {
	Enabled  bool
	CronExpr string
}

type jellyfinConfig struct {
	Url                                 string
	ApiKey                              string
	WatchedFilmFolders                  []string
	WatchedSeriesFolders                []string
	ObservedPeriodDays                  int
	IgnoreItemsAddedAfterLastNewsletter bool
}

type tmdbConfig struct {
	ApiKey string
}

type emailTemplateConfig struct {
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

type smtpConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	SenderName string
	TlsType    string
}

type dryRunConfig struct {
	Enabled            bool
	TestSTMPConnection bool
	OutputDirectory    string
	OutputFilename     string
	IncludeMetadata    bool
	SaveEmailData      bool
}

type Configuration struct {
	Debug           bool
	EmailRecipients []string
	Scheduler       schedulerConfig
	Jellyfin        jellyfinConfig
	Tmdb            tmdbConfig
	EmailTemplate   emailTemplateConfig
	SMTP            smtpConfig
	DryRun          dryRunConfig
}
