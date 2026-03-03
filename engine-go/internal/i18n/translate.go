package i18n

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localizer struct {
	*i18n.Localizer
}

func NewLocalizer(lang string) *Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	var translationFS embed.FS
	entries, _ := translationFS.ReadDir("internal/i18n")
	for _, e := range entries {
		bundle.LoadMessageFileFS(translationFS, "internal/i18n/"+e.Name())
	}

	return &Localizer{
		i18n.NewLocalizer(bundle, lang, "en"),
	}
}

func (l *Localizer) Localize(keyName string) string {
	translation, err := l.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID: keyName,
	})

	if err != nil {
		return "{" + keyName + "}"
	}

	return translation
}
