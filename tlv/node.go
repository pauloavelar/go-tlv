package tlv

import (
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/pauloavelar/go-tlv/tlv/internal/sizes"
	"github.com/pauloavelar/go-tlv/tlv/internal/utils"
)

// Node structure used to represent a decoded TLV message.
type Node struct {
	Tag Tag
	// Length of value part in bytes.
	Length Length
	// The underlying bytes with value part only.
	Value []byte
	// The underlying bytes with tag, length and value parts.
	Raw []byte

	decoder Decoder
}

// Tag node identifier composed by 1 to 8 bytes (uint64).
type Tag uint64

// Length value size in bytes.
type Length uint64

// Base64 converts the node bytes to base64.
func (n *Node) Base64() string {
	return base64.StdEncoding.EncodeToString(n.Raw)
}

// GetNodes parses the value as decoded TLV nodes.
func (n *Node) GetNodes() (Nodes, error) {
	return n.getSafeDecoder().DecodeAll(n.Value)
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

func (n *Node) getSafeDecoder() Decoder {
	if n.decoder != nil {
		return n.decoder
	}

	return stdDecoder
}

func (n *Node) getByteOrder() binary.ByteOrder {
	return n.getSafeDecoder().GetByteOrder()
}
