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
	count := 1 // default
	if len(pluralCount) > 0 {
		count = pluralCount[0]
	}
	translation, err := l.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID:   keyName,
		PluralCount: count,
	})

	if err != nil {
		return "{" + keyName + "}"
	}

	return translation
}
