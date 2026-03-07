package i18n

import (
	"embed"
	"errors"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localizer struct {
	*i18n.Localizer
}

//go:embed *.toml
var translationFS embed.FS

func NewLocalizer(lang string) (*Localizer, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	entries, err := translationFS.ReadDir(".")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		bundle.LoadMessageFileFS(translationFS, e.Name())
	}

	supportedLangs := []string{}
	isLangAvailable := false
	for _, t := range bundle.LanguageTags() {
		supportedLangs = append(supportedLangs, t.String())
		if t.String() == lang {
			isLangAvailable = true
		}
	}

	if !isLangAvailable {
		err = errors.New(lang + " is not a supported language. Supported languages are " + strings.Join(supportedLangs, ", "))
	}

	return &Localizer{
		i18n.NewLocalizer(bundle, lang, "en"),
	}, err
}

func (l *Localizer) Localize(keyName string, pluralCount ...int) string {
	localizeConfig := &i18n.LocalizeConfig{
		MessageID: keyName,
	}

	if len(pluralCount) > 0 {
		localizeConfig.PluralCount = pluralCount[0]
	}

	translation, _ := l.Localizer.Localize(localizeConfig)

	// i18n does its best to return a string. If there is a risk its grammatically incorrect, it will return an err and return a fallback string.
	// example https://github.com/nicksnyder/go-i18n/issues/241
	// For Jellyfin-Newsletter we prefer to have something incorrect but to have a string anyway instead of the placeholder.
	if translation == "" {
		return "{" + keyName + "}"
	}

	return translation
}
