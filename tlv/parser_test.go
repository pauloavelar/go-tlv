package tlv

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type panicReader struct{}

func (*panicReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("forcing reader error")
}

func TestParseReader(t *testing.T) {
	reader := bytes.NewReader(data)

	nodes, err := ParseReader(reader)

	require.Nil(t, err)
	require.Equal(t, 1, len(nodes))
	require.Equal(t, Tag(0x0001), nodes[0].Tag)
	require.Equal(t, uint16(71), nodes[0].Length)
	require.Equal(t, 71, len(nodes[0].Value))
}

func TestParseReader_WhenTheReaderFails(t *testing.T) {
	reader := new(panicReader)

	nodes, err := ParseReader(reader)

	require.NotNil(t, err)
	require.Nil(t, nodes)
}
