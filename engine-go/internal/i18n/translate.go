package i18n

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewLocalizer(lang string) *i18n.Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	var translationFS embed.FS
	entries, _ := translationFS.ReadDir("internal/i18n")
	for _, e := range entries {
		bundle.LoadMessageFileFS(translationFS, "internal/i18n/"+e.Name())
	}

	return i18n.NewLocalizer(bundle, lang, "en")
}
