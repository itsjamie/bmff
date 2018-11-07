package bmff

import (
	"encoding/binary"
)

type HintMediaHeader struct {
	*fullbox
	MaxPDUSize uint16
	AvgPDUSize uint16
	MaxBitrate uint32
	AvgBitrate uint32
}

func (b *HintMediaHeader) parse() error {
	b.MaxPDUSize = binary.BigEndian.Uint16(b.raw[0:2])
	b.AvgPDUSize = binary.BigEndian.Uint16(b.raw[2:4])
	b.MaxBitrate = binary.BigEndian.Uint32(b.raw[4:8])
	b.AvgBitrate = binary.BigEndian.Uint32(b.raw[8:12])

	b.raw = nil
	return nil
}
