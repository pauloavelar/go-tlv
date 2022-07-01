package tlv

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustCreateDecoder_WhenTheSizesAreInvalid(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()

	_ = MustCreateDecoder(0, 0, binary.BigEndian)
}

func TestCreateDecoder_WhenTheTagSizeIsTooSmall(t *testing.T) {
	d, err := CreateDecoder(0, 2, binary.BigEndian)

	require.Nil(t, d)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "tag size")
}

func TestCreateDecoder_WhenTheTagSizeIsTooBig(t *testing.T) {
	d, err := CreateDecoder(10, 2, binary.BigEndian)

	require.Nil(t, d)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "tag size")
}

func TestCreateDecoder_WhenTheLengthSizeIsTooSmall(t *testing.T) {
	d, err := CreateDecoder(2, 0, binary.BigEndian)

	require.Nil(t, d)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "length size")
}

func TestCreateDecoder_WhenTheLengthSizeIsTooBig(t *testing.T) {
	d, err := CreateDecoder(2, 10, binary.BigEndian)

	require.Nil(t, d)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "length size")
}
