package tlv

import (
	"encoding/binary"
	"io"

	"github.com/pauloavelar/go-tlv/tlv/internal/sizes"
)

// stdBinParser is the default value parser
var stdBinParser = binary.BigEndian

// stdDecoder uses 2 bytes for tags and lengths and parses them as big endian.
var stdDecoder = MustCreateDecoder(sizes.Uint16, sizes.Uint16, stdBinParser)

// DecodeReader decodes the whole reader as a list of TLV nodes.
func DecodeReader(reader io.Reader) (Nodes, error) {
	return stdDecoder.DecodeReader(reader)
}

// DecodeBytes decodes a byte array as a list of TLV nodes.
func DecodeBytes(data []byte) (Nodes, error) {
	return stdDecoder.DecodeBytes(data)
}

// DecodeSingle decodes a byte array as a single TLV node.
func DecodeSingle(data []byte) (res Node, read uint64, err error) {
	return stdDecoder.DecodeSingle(data)
}

// NewNode creates a node with the standard decoder configuration
func NewNode(tag Tag, value []byte) Node {
	return stdDecoder.NewNode(tag, value)
}
