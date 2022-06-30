package errors

import "fmt"

func NewInvalidSizeError(field string, value, min, max uint8) error {
	return fmt.Errorf("invalid %s size: %d (must be between %d and %d)", field, value, min, max)
}

func NewLengthMismatchError(expected uint64, message []byte, minNodeSize uint8) error {
	return fmt.Errorf(
		"value length mismatch, expected %d bytes but only %d bytes are available, data may be corrupted",
		expected, len(message)-int(minNodeSize),
	)
}

func NewMessageTooShortError(message []byte) error {
	return fmt.Errorf("message is too short (%d bytes), data may be corrupted", len(message))
}
