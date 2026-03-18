package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEmailAddressFromFriendlyName(t *testing.T) {
	tests := []struct {
		name            string
		address         string
		expectedAddress string
		expectError     bool
	}{
		{
			name:            "Complete friendly address",
			address:         "Test <test@domain.com>",
			expectedAddress: "test@domain.com",
			expectError:     false,
		},
		{
			name:            "Addresses with <> only",
			address:         "<test@domain.com>",
			expectedAddress: "test@domain.com",
			expectError:     false,
		},
		{
			name:            "Already cleaned address",
			address:         "test@domain.com",
			expectedAddress: "test@domain.com",
			expectError:     false,
		},
		{
			name:            "Invalid address",
			address:         "Test <@domain>",
			expectedAddress: "",
			expectError:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanedAddr, err := getEmailAddressFromFriendlyName(test.address)
			if test.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedAddress, cleanedAddr)
			}
		})
	}
}
