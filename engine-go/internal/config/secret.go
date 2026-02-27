package config

import "encoding/xml"

// This is a dedicated type to handle secrets loaded in struct
// The main goal is to avoid JSON serialization or logging of secrets

type Secret string

func (s Secret) String() string { return "[REDACTED]" }

func (s Secret) MarshalJSON() ([]byte, error) {
	return []byte(`"[REDACTED]"`), nil
}

func (s Secret) MarshalYAML() (any, error) { return "[REDACTED]", nil }

func (s Secret) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement("[REDACTED]", start)
}
func (s Secret) MarshalText() ([]byte, error) { return []byte("[REDACTED]"), nil }
