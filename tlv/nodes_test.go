package tlv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodes_HasTag_WhenTheTagIsPresent(t *testing.T) {
	tag := Tag(0x1234)
	nodes := Nodes{
		Node{Tag: tag, Value: []byte{}},
		Node{Tag: 0x4321, Value: []byte{}},
	}

	require.True(t, nodes.HasTag(tag))
}

func TestNodes_HasTag_WhenTheTagIsNotPresent(t *testing.T) {
	nodes := Nodes{
		Node{Tag: 0x1, Value: []byte{}},
		Node{Tag: 0x2, Value: []byte{}},
	}

	require.False(t, nodes.HasTag(0x3))
}
