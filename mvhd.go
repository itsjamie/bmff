package bmff

import (
	"encoding/binary"

	"github.com/itsjamie/bmff/internal/fixed"
	"github.com/pkg/errors"
)

type MovieHeader struct {
	*box
	version          uint8
	flags            [3]byte
	CreationTime     uint64
	ModificationTime uint64
	TimeScale        uint32
	Duration         uint64
	NextTrackID      uint32
	Rate             fixed.Uint16_16
	Volume           fixed.Uint8_8
	Reserved         []byte
	Matrix           [9]int32
	Predefined       []byte
}

func (b *MovieHeader) parse() error {
	b.version = b.raw[0]
	b.flags = [3]byte{b.raw[1], b.raw[2], b.raw[3]}
	var offset int
	if b.version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(b.raw[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(b.raw[8:12]))
		b.TimeScale = binary.BigEndian.Uint32(b.raw[12:16])
		b.Duration = uint64(binary.BigEndian.Uint32(b.raw[16:20]))
		offset = 20
	} else if b.version == 1 {
		b.CreationTime = binary.BigEndian.Uint64(b.raw[4:12])
		b.ModificationTime = binary.BigEndian.Uint64(b.raw[12:20])
		b.TimeScale = binary.BigEndian.Uint32(b.raw[20:24])
		b.Duration = binary.BigEndian.Uint64(b.raw[24:32])
		offset = 32
	}

	if err := b.Rate.UnmarshalBinary(b.raw[offset : offset+4]); err != nil {
		return errors.Wrap(err, "failed to get rate")
	}

	if err := b.Volume.UnmarshalBinary(b.raw[offset+4 : offset+6]); err != nil {
		return errors.Wrap(err, "failed to get volume")
	}

	b.Reserved = b.raw[offset+6 : offset+16]
	offset += 16
	for i := 0; i < 9; i++ {
		b.Matrix[i] = int32(binary.BigEndian.Uint32(b.raw[offset+i : offset+i+4]))
	}
	offset += 36

	b.Predefined = b.raw[offset : offset+24]
	b.NextTrackID = binary.BigEndian.Uint32(b.raw[offset+24 : offset+28])
	return nil
}
