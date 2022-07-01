package utils

import (
	"encoding/binary"

	"github.com/pauloavelar/go-tlv/tlv/internal/sizes"
)

func GetPadding(typeSize, valueSize int) []byte {
	padSize := MaxInt(0, typeSize-valueSize)
	return make([]byte, padSize, typeSize)
}

func GetPaddedUint64(byteOrder binary.ByteOrder, data []byte) uint64 {
	padding := GetPadding(sizes.Uint64, len(data))
	return byteOrder.Uint64(append(padding, data...))
}
