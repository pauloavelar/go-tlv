package tlv

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pauloavelar/go-tlv/tlv/internal/sizes"
)

const (
	tagSize     = sizes.Uint16
	lengthSize  = sizes.Uint16
	minNodeSize = tagSize + lengthSize
)

// parser bit parsing defaults to BigEndian
var parser = binary.BigEndian

// Parse a Reader to TLV nodes
func ParseReader(reader io.Reader) (Nodes, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return ParseBytes(data)
}

// ParseBytes parses a byte array to TLV nodes
func ParseBytes(data []byte) (Nodes, error) {
	node, read, err := ParseSingle(data)
	if err != nil {
		return nil, err
	}

	if len(data) == read {
		return Nodes{node}, nil
	}

	next, err := ParseBytes(data[read:])
	if err != nil {
		return nil, err
	}

	return append(Nodes{node}, next...), nil
}

// ParseSingle parses a byte array to a single TLV node
func ParseSingle(data []byte) (res Node, read int, err error) {
	if len(data) < minNodeSize {
		return res, 0, fmt.Errorf("message is too short (%d bytes), data may be corrupted", len(data))
	}

	tag := parser.Uint16(data[:tagSize])
	length := parser.Uint16(data[tagSize:minNodeSize])
	messageLength := minNodeSize + int(length)

	if len(data) < messageLength {
		return res, 0, fmt.Errorf(
			"value length mismatch, expected %d bytes but only %d bytes are available, data may be corrupted",
			length, len(data)-minNodeSize,
		)
	}

	node := Node{
		Tag:    Tag(tag),
		Length: length,
		Value:  data[minNodeSize:messageLength],
		Raw:    data[:messageLength],
	}

	return node, messageLength, nil
}
