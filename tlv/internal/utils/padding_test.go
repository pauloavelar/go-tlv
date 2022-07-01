package utils

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPadding(t *testing.T) {
	padding := GetPadding(4, 2)

	require.Equal(t, []byte{0, 0}, padding)
}

func TestGetPaddedUint64(t *testing.T) {
	value := GetPaddedUint64(binary.BigEndian, []byte{0x01, 0x23, 0x45})

	require.EqualValues(t, 0x12345, value)
}
