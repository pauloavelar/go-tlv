package utils

func GetPadding(typeSize, valueSize int) []byte {
	padSize := MaxInt(0, typeSize-valueSize)
	return make([]byte, padSize, typeSize)
}
