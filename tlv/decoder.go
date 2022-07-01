package tlv

import (
	"encoding/binary"
	"io"

	"github.com/pauloavelar/go-tlv/tlv/internal/errors"
	"github.com/pauloavelar/go-tlv/tlv/internal/utils"
)

// Decoder is a TLV decoder with custom configuration.
type Decoder interface {
	// DecodeReader decodes the whole reader to a list of TLV nodes
	DecodeReader(reader io.Reader) (Nodes, error)
	// DecodeBytes decodes a byte array to a list of TLV nodes
	DecodeBytes(data []byte) (Nodes, error)
	// DecodeSingle decodes a byte array to a single TLV node
	DecodeSingle(data []byte) (res Node, read uint64, err error)
}

type decoder struct {
	tagSize      uint8
	lengthSize   uint8
	minNodeSize  uint8
	binaryParser binary.ByteOrder
}

const (
	minTagSize = 1 // 2^1 = 2
	maxTagSize = 8 // 2^8 = 256
	minLenSize = 1 // 2^1 = 2
	maxLenSize = 8 // 2^8 = 256
)

// MustCreateDecoder creates a decoder using custom configuration or panics in case of any errors.
func MustCreateDecoder(tagSize, lengthSize uint8, byteOrder binary.ByteOrder) Decoder {
	res, err := CreateDecoder(tagSize, lengthSize, byteOrder)
	if err != nil {
		panic(err)
	}

	return res
}

// CreateDecoder creates a decoder using custom configuration.
// Hint: `tagSize` and `lengthSize` must be numbers between 1 and 8.
func CreateDecoder(tagSize, lengthSize uint8, byteOrder binary.ByteOrder) (Decoder, error) {
	if tagSize < minTagSize || tagSize > maxTagSize {
		return nil, errors.NewInvalidSizeError("tag", tagSize, minTagSize, maxTagSize)
	}

	if lengthSize < minLenSize || lengthSize > maxLenSize {
		return nil, errors.NewInvalidSizeError("length", lengthSize, minLenSize, maxLenSize)
	}

	res := &decoder{
		tagSize:      tagSize,
		lengthSize:   lengthSize,
		minNodeSize:  tagSize + lengthSize,
		binaryParser: byteOrder,
	}

	return res, nil
}

// DecodeReader decodes the full contents of a Reader as TLV nodes.
// Note: the current implementation loads the entire Reader data into memory.
func (p *decoder) DecodeReader(reader io.Reader) (Nodes, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return p.DecodeBytes(data)
}

// DecodeBytes decodes a byte array as TLV nodes.
func (p *decoder) DecodeBytes(data []byte) (Nodes, error) {
	node, read, err := p.DecodeSingle(data)
	if err != nil {
		return nil, err
	}

	if uint64(len(data)) == read {
		return Nodes{node}, nil
	}

	next, err := p.DecodeBytes(data[read:])
	if err != nil {
		return nil, err
	}

	return append(Nodes{node}, next...), nil
}

// DecodeSingle decodes a byte array as a single TLV node.
func (p *decoder) DecodeSingle(data []byte) (res Node, read uint64, err error) {
	if len(data) < int(p.minNodeSize) {
		return res, 0, errors.NewMessageTooShortError(data)
	}

	tag := utils.GetPaddedUint64(p.binaryParser, data[:p.tagSize])
	length := utils.GetPaddedUint64(p.binaryParser, data[p.tagSize:p.minNodeSize])
	messageLength := uint64(p.minNodeSize) + length

	if len(data) < int(messageLength) {
		return res, 0, errors.NewLengthMismatchError(length, data, p.minNodeSize)
	}

	node := Node{
		Tag:       Tag(tag),
		Length:    Length(length),
		Value:     data[p.minNodeSize:messageLength],
		Raw:       data[:messageLength],
		decoder:   p,
		binParser: p.binaryParser,
	}

	return node, messageLength, nil
}
