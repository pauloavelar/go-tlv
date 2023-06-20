package tlv

import (
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/pauloavelar/go-tlv/tlv/internal/errors"
	"github.com/pauloavelar/go-tlv/tlv/internal/sizes"
	"github.com/pauloavelar/go-tlv/tlv/internal/utils"
)

// Node structure used to represent a decoded TLV message.
type Node struct {
	Tag    Tag
	Length Length
	Value  []byte
	Raw    []byte

	decoder Decoder
}

// Tag node identifier composed by 1 to 8 bytes (uint64).
type Tag uint64

// Length value size in bytes.
type Length uint64

// String converts the node bytes to base64.
func (n *Node) String() string {
	return base64.StdEncoding.EncodeToString(n.Raw)
}

// GetNodes parses the value as decoded TLV nodes.
func (n *Node) GetNodes() (Nodes, error) {
	return n.getSafeDecoder().DecodeBytes(n.Value)
}

// GetBool parses the value as boolean if it has enough bytes.
func (n *Node) GetBool() (res, ok bool) {
	if len(n.Value) < sizes.Bool {
		return false, false
	}

	return n.Value[0] != 0, true
}

// GetPaddedBool parses the value as boolean regardless of its size.
func (n *Node) GetPaddedBool() bool {
	res, _ := n.GetBool()
	return res
}

// GetString parses the value as UTF-8 text.
func (n *Node) GetString() string {
	return string(n.Value)
}

// GetDate parses the value as date if it has enough bytes.
func (n *Node) GetDate() (res time.Time, ok bool) {
	if len(n.Value) == 0 {
		return res, false
	}

	epoch := n.GetPaddedUint64()
	return time.Unix(int64(epoch), 0).UTC(), true
}

// GetUint8 parses the value as uint8.
func (n *Node) GetUint8() (res uint8, ok bool) {
	if len(n.Value) < sizes.Uint8 {
		return 0, false
	}

	return n.Value[0], true
}

// GetPaddedUint8 parses the value as uint8 regardless of size.
func (n *Node) GetPaddedUint8() uint8 {
	if len(n.Value) < sizes.Uint8 {
		return 0
	}

	return n.Value[0]
}

// GetUint16 parses the value as uint16 if it has enough bytes.
func (n *Node) GetUint16() (res uint16, ok bool) {
	if len(n.Value) < sizes.Uint16 {
		return 0, false
	}

	return n.getByteOrder().Uint16(n.Value), true
}

// GetPaddedUint16 parses the value as uint16 regardless of size.
func (n *Node) GetPaddedUint16() uint16 {
	padding := utils.GetPadding(sizes.Uint16, len(n.Value))

	return n.getByteOrder().Uint16(append(padding, n.Value...))
}

// GetUint32 parses the value as uint32 if it has enough bytes.
func (n *Node) GetUint32() (res uint32, exists bool) {
	if len(n.Value) < sizes.Uint32 {
		return 0, false
	}

	return n.getByteOrder().Uint32(n.Value), true
}

// GetPaddedUint32 parses the value as uint32 regardless of size.
func (n *Node) GetPaddedUint32() uint32 {
	padding := utils.GetPadding(sizes.Uint32, len(n.Value))

	return n.getByteOrder().Uint32(append(padding, n.Value...))
}

// GetUint64 parses the value as uint64 if it has enough bytes.
func (n *Node) GetUint64() (res uint64, ok bool) {
	if len(n.Value) < sizes.Uint64 {
		return 0, false
	}

	return n.getByteOrder().Uint64(n.Value), true
}

// GetPaddedUint64 parses the value as uint64 regardless of size.
func (n *Node) GetPaddedUint64() uint64 {
	padding := utils.GetPadding(sizes.Uint64, len(n.Value))

	return n.getByteOrder().Uint64(append(padding, n.Value...))
}

// GetVariantArray parses the value as an array. All nodes in the array have individual tag.
func (n *Node) GetVariantArray() (Nodes, error) {
	lengthSize := uint64(n.decoder.GetLengthSize())
	if len(n.Value) < int(lengthSize) {
		return nil, errors.NewMessageTooShortError(n.Value)
	}
	arrayLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[:lengthSize])
	nodes := make([]Node, 0, arrayLength)
	offset := lengthSize

	for i := uint64(0); i < arrayLength; i++ {
		node, read, err := n.decoder.DecodeSingle(n.Value[offset:])
		if err != nil {
			return nil, err
		}
		offset += read
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// GetArray parses the value as an array. All nodes in the array have the same tag.
func (n *Node) GetArray() (Nodes, error) {
	tagSize := uint64(n.decoder.GetTagSize())
	lengthSize := uint64(n.decoder.GetLengthSize())
	if len(n.Value) < int(lengthSize+tagSize) {
		return nil, errors.NewMessageTooShortError(n.Value)
	}
	arrayTag := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[:tagSize])
	arrayLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[tagSize:tagSize+lengthSize])
	nodes := make([]Node, 0, arrayLength)
	offset := tagSize + lengthSize

	for i := uint64(0); i < arrayLength; i++ {
		itemLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[offset:offset+lengthSize])
		offset += lengthSize
		nodes = append(nodes, Node{
			Tag:    Tag(arrayTag),
			Length: Length(itemLength),
			Value:  n.Value[offset : offset+itemLength],
			// All the array items use the same tag, so we do not provide raw bytes.
			Raw: nil,
		})
		offset += itemLength
	}
	return nodes, nil
}

// GetStringMap parses the value as an key-value pair, the key type is string, and
// each value have individual tag.
func (n *Node) GetVariantStringMap() (map[string]*Node, error) {
	lengthSize := uint64(n.decoder.GetLengthSize())
	if len(n.Value) < int(lengthSize) {
		return nil, errors.NewMessageTooShortError(n.Value)
	}
	mapLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[:lengthSize])
	nodes := make(map[string]*Node, mapLength)
	offset := lengthSize

	for i := uint64(0); i < mapLength; i++ {
		labelLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[offset:offset+lengthSize])
		offset += lengthSize
		labelString := string(n.Value[offset : offset+labelLength])
		offset += labelLength

		node, read, err := n.decoder.DecodeSingle(n.Value[offset:])
		if err != nil {
			return nil, err
		}
		offset += read
		nodes[labelString] = &node
	}
	return nodes, nil
}

// GetStringMap parses the value as an key-value pair, the key type is string, and
// each value have the same tag.
func (n *Node) GetStringMap() (map[string]*Node, error) {
	tagSize := uint64(n.decoder.GetTagSize())
	lengthSize := uint64(n.decoder.GetLengthSize())
	if len(n.Value) < int(lengthSize+tagSize) {
		return nil, errors.NewMessageTooShortError(n.Value)
	}
	mapTag := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[:tagSize])
	mapLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[tagSize:tagSize+lengthSize])
	nodes := make(map[string]*Node, mapLength)
	offset := tagSize + lengthSize

	for i := uint64(0); i < mapLength; i++ {
		labelLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[offset:offset+lengthSize])
		offset += lengthSize
		labelString := string(n.Value[offset : offset+labelLength])
		offset += labelLength

		itemLength := utils.GetPaddedUint64(n.decoder.GetByteOrder(), n.Value[offset:offset+lengthSize])
		offset += lengthSize
		nodes[labelString] = &Node{
			Tag:    Tag(mapTag),
			Length: Length(itemLength),
			Value:  n.Value[offset : offset+itemLength],
			// All the map items use the same tag, so we do not provide raw bytes.
			Raw: nil,
		}
		offset += itemLength
	}
	return nodes, nil
}

func (n *Node) getSafeDecoder() Decoder {
	if n.decoder != nil {
		return n.decoder
	}

	return stdDecoder
}

func (n *Node) getByteOrder() binary.ByteOrder {
	return n.getSafeDecoder().GetByteOrder()
}
