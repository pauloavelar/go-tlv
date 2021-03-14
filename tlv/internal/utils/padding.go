package utils

func GetPadding(typeSize int, valueSize int) []byte {
	padSize := MaxInt(0, typeSize-valueSize)
	return make([]byte, padSize, typeSize)
}
