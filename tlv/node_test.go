package tlv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_String(t *testing.T) {
	node := Node{Raw: []byte{0x01, 0x02, 0x03, 0x04, 0x05}}

	require.Equal(t, "AQIDBAU=", node.String())
}

func TestNode_GetBool_WhenValueIsTrue(t *testing.T) {
	node := Node{Value: []byte{0x01}}

	res, ok := node.GetBool()

	require.True(t, ok)
	require.True(t, res)
}

func TestNode_GetBool_WhenValueIsFalse(t *testing.T) {
	node := Node{Value: []byte{0x00}}

	res, ok := node.GetBool()

	require.True(t, ok)
	require.False(t, res)
}

func TestNode_GetBool_WhenValueIsEmpty(t *testing.T) {
	node := Node{Value: []byte{}}

	res, ok := node.GetBool()

	require.False(t, ok)
	require.False(t, res)
}

func TestNode_GetPaddedBool(t *testing.T) {
	require.True(t, Node{Value: []byte{0x01}}.GetPaddedBool())
	require.False(t, Node{Value: []byte{0x00}}.GetPaddedBool())
	require.False(t, Node{Value: []byte{}}.GetPaddedBool())
	require.False(t, Node{Value: nil}.GetPaddedBool())
}
