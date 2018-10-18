package bmff

import (
	"encoding/binary"

	"github.com/itsjamie/bmff/internal/fixed"
)

type TrackHeader struct {
	*box
	version          uint8
	flags            [3]byte
	CreationTime     uint64
	ModificationTime uint64
	TrackID          uint32
	Duration         uint64
	Layer            int16
	AlternateGroup   int16
	Volume           int16
	Matrix           [9]int32
	Width            fixed.Uint16_16
	Height           fixed.Uint16_16
}

func (b *TrackHeader) parse() error {
	b.version = b.raw[0]
	b.flags = [3]byte{b.raw[1], b.raw[2], b.raw[3]}
	var offset int
	if b.version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(b.raw[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(b.raw[8:12]))
		b.TrackID = binary.BigEndian.Uint32(b.raw[12:16])
		// 16:20 reserved
		b.Duration = uint64(binary.BigEndian.Uint32(b.raw[20:24]))
		offset = 24
	} else if b.version == 1 {
		b.CreationTime = uint64(binary.BigEndian.Uint64(b.raw[4:12]))
		b.ModificationTime = uint64(binary.BigEndian.Uint64(b.raw[12:20]))
		b.TrackID = binary.BigEndian.Uint32(b.raw[20:24])
		// 24:28 reserved
		b.Duration = uint64(binary.BigEndian.Uint64(b.raw[28:36]))
		offset = 36
	}
	offset += 8 // reserved bytes
	b.Layer = int16(binary.BigEndian.Uint16(b.raw[offset : offset+2]))
	b.AlternateGroup = int16(binary.BigEndian.Uint16(b.raw[offset+2 : offset+4]))
	b.Volume = int16(binary.BigEndian.Uint16(b.raw[offset+4 : offset+6]))
	offset += 8 // previous bytes + 2 reserved

	for i := 0; i < 9; i++ {
		b.Matrix[i] = int32(binary.BigEndian.Uint32(b.raw[offset+i : offset+i+4]))
	}
	offset += 36
	b.Width = fixed.Uint16_16(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
	b.Height = fixed.Uint16_16(binary.BigEndian.Uint32(b.raw[offset+4 : offset+8]))
	return nil
}
