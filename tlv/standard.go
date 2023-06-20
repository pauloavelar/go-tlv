package tlv

import (
	"encoding/binary"
	"io"

	"github.com/pauloavelar/go-tlv/tlv/internal/sizes"
)

// stdByteOrder is the default endianness for parsing numbers.
var stdByteOrder = binary.BigEndian

// stdDecoder uses 2 bytes for tags and lengths and parses them as big endian.
var stdDecoder = MustCreateDecoder(sizes.Uint16, sizes.Uint16, stdByteOrder)

// DecodeReader decodes the entire [io.Reader] data as a list of TLV nodes.
func DecodeReader(reader io.Reader) (Nodes, error) {
	return stdDecoder.DecodeReader(reader)
}

// DecodeAll decodes a byte array as a list of TLV [Nodes].
func DecodeAll(data []byte) (Nodes, error) {
	return stdDecoder.DecodeAll(data)
}

// DecodeSingle decodes a byte array as a single TLV [Node].
func DecodeSingle(data []byte) (res Node, read uint64, err error) {
	return stdDecoder.DecodeSingle(data)
}

// NewNode creates a [Node] with the default [Decoder] configuration.
func NewNode(tag Tag, value []byte) Node {
	return stdDecoder.NewNode(tag, value)
}
