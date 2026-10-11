package config

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecret_SafeString(t *testing.T) {
	tests := []struct {
		name string
		s    Secret
		want string
	}{
		{name: "non-empty secret", s: Secret("super-secret-value"), want: "super-secret-value"},
		{name: "empty secret", s: Secret(""), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.SafeString())
		})
	}
}

func TestSecret_String_Redacted(t *testing.T) {
	secret := Secret("super-secret-value")
	assert.Equal(t, "[REDACTED]", secret.String())
	assert.NotContains(t, secret.String(), string(secret))

	emptySecret := Secret("")
	assert.Equal(t, "[REDACTED]", emptySecret.String())
}

func TestSecret_MarshalJSON_Redacted(t *testing.T) {
	tests := []struct {
		name string
		s    Secret
	}{
		{name: "non-empty secret", s: Secret("super-secret-value")},
		{name: "empty secret", s: Secret("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.Marshal(tt.s)
			require.NoError(t, err)
			assert.Equal(t, `"[REDACTED]"`, string(out))
			if tt.s != "" {
				assert.NotContains(t, string(out), string(tt.s))
			}
		})
	}
}

func TestSecret_MarshalYAML_Redacted(t *testing.T) {
	tests := []struct {
		name string
		s    Secret
	}{
		{name: "non-empty secret", s: Secret("super-secret-value")},
		{name: "empty secret", s: Secret("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := tt.s.MarshalYAML()
			require.NoError(t, err)
			assert.Equal(t, "[REDACTED]", out)
		})
	}
}

func TestSecret_MarshalText_Redacted(t *testing.T) {
	tests := []struct {
		name string
		s    Secret
	}{
		{name: "non-empty secret", s: Secret("super-secret-value")},
		{name: "empty secret", s: Secret("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := tt.s.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, "[REDACTED]", string(out))
			if tt.s != "" {
				assert.NotContains(t, string(out), string(tt.s))
			}
		})
	}
}

func TestSecret_MarshalXML_Redacted(t *testing.T) {
	tests := []struct {
		name string
		s    Secret
	}{
		{name: "non-empty secret", s: Secret("super-secret-value")},
		{name: "empty secret", s: Secret("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := xml.NewEncoder(&buf)
			require.NoError(t, enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "secret"}}))
			require.NoError(t, tt.s.MarshalXML(enc, xml.StartElement{Name: xml.Name{Local: "secret"}}))
			require.NoError(t, enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "secret"}}))
			require.NoError(t, enc.Flush())

			assert.Contains(t, buf.String(), "[REDACTED]")
			if tt.s != "" {
				assert.NotContains(t, buf.String(), string(tt.s))
			}
		})
	}
}

func TestSecret_MarshalXML_FullDocument(t *testing.T) {
	type config struct {
		APIKey Secret `xml:"apiKey"`
	}

	cfg := config{APIKey: Secret("super-secret-value")}
	out, err := xml.Marshal(cfg)
	require.NoError(t, err)

	assert.Contains(t, string(out), "[REDACTED]")
	assert.NotContains(t, string(out), "super-secret-value")
}

func TestSecret_InStruct_MarshalJSON_Redacted(t *testing.T) {
	type cfg struct {
		APIKey Secret `json:"apiKey"`
	}

	c := cfg{APIKey: Secret("super-secret-value")}
	out, err := json.Marshal(c)
	require.NoError(t, err)

	assert.Contains(t, string(out), `"[REDACTED]"`)
	assert.NotContains(t, string(out), "super-secret-value")
}

func TestSecret_FmtRedaction(t *testing.T) {
	s := Secret("super-secret-value")
	formatted := strings.TrimSpace(s.String())
	assert.Equal(t, "[REDACTED]", formatted)
}
