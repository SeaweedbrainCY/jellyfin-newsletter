package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
