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
		_, err := bundle.LoadMessageFileFS(translationFS, e.Name())
		if err != nil {
			return nil, err
		}
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
		err = errors.New(
			lang + " is not a supported language. Supported languages are " + strings.Join(supportedLangs, ", "),
		)
	}

	return &Localizer{
		i18n.NewLocalizer(bundle, lang, "en"),
	}, err
}

func (l *Localizer) getLocalization(config *i18n.LocalizeConfig) string {
	translation, _ := l.Localizer.Localize(config)

	// i18n does its best to return a string. If there is a risk its grammatically incorrect, it will return an err and return a fallback string.
	// example https://github.com/nicksnyder/go-i18n/issues/241
	// For Jellyfin-Newsletter we prefer to have something incorrect but to have a string anyway instead of the placeholder.
	if translation == "" {
		return "{" + config.MessageID + "}"
	}

	return translation
}

func (l *Localizer) LocalizeWithTemplate(keyName string, templateData interface{}) string {
	localizeConfig := &i18n.LocalizeConfig{
		MessageID:    keyName,
		TemplateData: templateData,
	}
	return l.getLocalization(localizeConfig)
}

func (l *Localizer) LocalizeWithPlural(keyName string, pluralCount int) string {
	localizeConfig := &i18n.LocalizeConfig{
		MessageID:   keyName,
		PluralCount: pluralCount,
	}
	return l.getLocalization(localizeConfig)
}

func (l *Localizer) Localize(keyName string) string {
	localizeConfig := &i18n.LocalizeConfig{
		MessageID: keyName,
	}
	return l.getLocalization(localizeConfig)
}
