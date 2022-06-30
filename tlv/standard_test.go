package tlv

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type failingReader struct{}

func (*failingReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("forcing reader error")
}

func TestDecodeReader(t *testing.T) {
	reader := bytes.NewReader(data)

	nodes, err := DecodeReader(reader)

	require.Nil(t, err)
	require.Equal(t, 1, len(nodes))
	require.Equal(t, Tag(0x0001), nodes[0].Tag)
	require.Equal(t, Length(71), nodes[0].Length)
	require.Equal(t, 71, len(nodes[0].Value))
}

func TestDecodeReader_WhenTheReaderFails(t *testing.T) {
	reader := new(failingReader)

	nodes, err := DecodeReader(reader)

	require.NotNil(t, err)
	require.Nil(t, nodes)
}

func TestDecodeSingle_WhenTheDataIsCorrupted(t *testing.T) {
	corrupted := data[:len(data)-5]

	node, read, err := DecodeSingle(corrupted)

	require.NotNil(t, err)
	require.Zero(t, read)
	require.Empty(t, node)
}

func TestDecodeBytes_WhenTheDataIsCorrupted(t *testing.T) {
	corrupted := make([]byte, 0, len(data)*2-5)
	corrupted = append(corrupted, data...)
	corrupted = append(corrupted, data[:len(data)-5]...)

	node, err := DecodeBytes(corrupted)

	require.NotNil(t, err)
	require.Nil(t, node)
}
