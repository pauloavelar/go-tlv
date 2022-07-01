package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInvalidSizeError(t *testing.T) {
	err := NewInvalidSizeError("field", 1, 2, 3)

	require.NotNil(t, err)
	require.Equal(t, "invalid field size: 1 (must be between 2 and 3)", err.Error())
}

func TestNewLengthMismatchError(t *testing.T) {
	expected := "value length mismatch, expected 5 bytes but only 4 bytes are available, data may be corrupted"

	err := NewLengthMismatchError(5, []byte("12345678"), 4)

	require.NotNil(t, err)
	require.Equal(t, expected, err.Error())
}

func TestNewMessageTooShortError(t *testing.T) {
	err := NewMessageTooShortError([]byte("123"))

	require.NotNil(t, err)
	require.Equal(t, "message is too short (3 bytes), data may be corrupted", err.Error())
}
