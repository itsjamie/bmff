package mp4bits

// Int2 represents a intemediary definition for getting 2-bit unsigned ints
type Int2 byte

// First returns a 2-bit integer in the first pair of bits
func (i Int2) First() uint8 {
	return uint8(i & 0b11000000 >> 6)
}

// Second returns a 2-bit integer in the second pair of bits
func (i Int2) Second() uint8 {
	return uint8(i & 0b00110000 >> 4)
}

// Third returns a 2-bit integer in the third pair of bits
func (i Int2) Third() uint8 {
	return uint8(i & 0b00001100 >> 2)
}

// Fourth returns a 2-bit integer in the fourth pair of bits
func (i Int2) Fourth() uint8 {
	return uint8(i & 0b00000011)
}
