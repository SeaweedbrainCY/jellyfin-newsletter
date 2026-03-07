package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLocalizerWithValidLang(t *testing.T) {
	tests := []struct {
		name string
		lang string
	}{
		{
			name: "Test Get Localizer for fr lang",
			lang: "fr",
		},
		{
			name: "Test Get Localizer for en lang",
			lang: "en",
		},
		{
			name: "Test Get Localizer for es lang",
			lang: "es",
		},
		{
			name: "Test Get Localizer for fi lang",
			lang: "fi",
		},
		{
			name: "Test Get Localizer for he lang",
			lang: "he",
		},
		{
			name: "Test Get Localizer for de lang",
			lang: "de",
		},
		{
			name: "Test Get Localizer for it lang",
			lang: "it",
		},
		{
			name: "Test Get Localizer for pt lang",
			lang: "pt",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewLocalizer(test.lang)
			assert.NoError(t, err)
		})
	}
}

func TestGetLocalizerWithUnknownLang(t *testing.T) {
	_, err := NewLocalizer("zz")
	assert.ErrorContains(t, err, "zz is not a supported language. Supported languages are")
}

func TestLocalizeValidStrings(t *testing.T) {
	key := "currently_available"
	tests := []struct {
		name           string
		lang           string
		expectedString string
	}{
		{
			name:           "Test Get Localization of currently_available in fr",
			lang:           "fr",
			expectedString: "Actuellement disponible sur Jellyfin :",
		},
		{
			name:           "Test Get Localization of currently_available in en",
			lang:           "en",
			expectedString: "Currently available in Jellyfin:",
		},
		{
			name:           "Test Get Localization of currently_available in es",
			lang:           "es",
			expectedString: "Disponible actualmente en Jellyfin:",
		},
		{
			name:           "Test Get Localization of currently_available in fi",
			lang:           "fi",
			expectedString: "Tällä hetkellä saatavilla Jellyfinissä:",
		},
		{
			name:           "Test Get Localization of currently_available in he",
			lang:           "he",
			expectedString: "זמין כעת בג'ליפין:\\u200f",
		},
		{
			name:           "Test Get Localization of currently_available in de",
			lang:           "de",
			expectedString: "Derzeit verfügbar in Jellyfin:",
		},
		{
			name:           "Test Get Localization of currently_available in it",
			lang:           "it",
			expectedString: "Attualmente disponibile su Jellyfin:",
		},
		{
			name:           "Test Get Localization of currently_available in pt",
			lang:           "pt",
			expectedString: "Atualmente disponível no Jellyfin:",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l, err := NewLocalizer(test.lang)
			require.NoError(t, err)
			localizedStr := l.Localize(key)
			assert.Equal(t, test.expectedString, localizedStr)
		})
	}
}

func TestLocalizeWithNotLocalizedString(t *testing.T) {
	// We create a Localizer for a lang that doesn't exist. Purposely ignore the error which will tell the lang doesn't exist. This way we are sure that the key we want to localize doesn't exist in our imaginary lang, but exists in english, the fallback.

	l, _ := NewLocalizer("zz")
	localizedStr := l.Localize("currently_available")
	assert.Equal(t, "Currently available in Jellyfin:", localizedStr)
}

func TestLocalizeUnknownKey(t *testing.T) {
	l, _ := NewLocalizer("fr")
	localizedStr := l.Localize("thisKeyDoesntExist")
	assert.Equal(t, "{thisKeyDoesntExist}", localizedStr)
}
