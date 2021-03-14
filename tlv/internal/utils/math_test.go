package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaxInt(t *testing.T) {
	require.Equal(t, 9, MaxInt(1, 9))
	require.Equal(t, 8, MaxInt(8, 2))
}
