package tlv

import (
	"encoding/binary"
	"io"

	"github.com/pauloavelar/go-tlv/tlv/internal/sizes"
)

// stdDecoder uses 2 bytes for tags and lengths and parses them as big endian.
var stdDecoder = MustCreateDecoder(sizes.Uint16, sizes.Uint16, binary.BigEndian)

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
