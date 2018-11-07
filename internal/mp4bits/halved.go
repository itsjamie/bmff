package mp4bits

// HalvedByte is a representation for two numbers stored in the same byte.
// Each take 4 bits, and are stored in BigEndian order.
type HalvedByte byte

// Low returns the last 4 bits of the byte
func (h HalvedByte) Low() uint8 {
	return uint8(h & 0x0F)
}

// High returns the first 4 bits of the byte
func (h HalvedByte) High() uint8 {
	return uint8(h >> 4)
}
