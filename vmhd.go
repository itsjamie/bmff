package bmff

import (
	"encoding/binary"
)

type VideoMediaHeader struct {
	*fullbox
	graphicsMode uint16
	opcolor      [3]uint16
}

func (b *VideoMediaHeader) parse() error {
	b.graphicsMode = binary.BigEndian.Uint16(b.raw[0:2])
	b.opcolor[0] = binary.BigEndian.Uint16(b.raw[2:4])
	b.opcolor[1] = binary.BigEndian.Uint16(b.raw[4:6])
	b.opcolor[2] = binary.BigEndian.Uint16(b.raw[6:8])
	b.raw = nil
	return nil
}
