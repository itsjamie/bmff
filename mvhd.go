package bmff

import (
	"encoding/binary"

	"github.com/itsjamie/bmff/internal/fixed"
	"github.com/pkg/errors"
)

type MovieHeader struct {
	*fullbox
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
	var offset int
	if b.version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(b.raw[0:4]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(b.raw[4:8]))
		b.TimeScale = binary.BigEndian.Uint32(b.raw[8:12])
		b.Duration = uint64(binary.BigEndian.Uint32(b.raw[12:16]))
		offset = 16
	} else if b.version == 1 {
		b.CreationTime = binary.BigEndian.Uint64(b.raw[0:8])
		b.ModificationTime = binary.BigEndian.Uint64(b.raw[8:16])
		b.TimeScale = binary.BigEndian.Uint32(b.raw[16:20])
		b.Duration = binary.BigEndian.Uint64(b.raw[20:28])
		offset = 28
	}

	if err := b.Rate.UnmarshalBinary(b.raw[offset : offset+4]); err != nil {
		return errors.Wrap(err, "failed to get rate")
	}
	offset += 4

	if err := b.Volume.UnmarshalBinary(b.raw[offset : offset+2]); err != nil {
		return errors.Wrap(err, "failed to get volume")
	}
	offset += 2

	b.Reserved = b.raw[offset : offset+10]
	offset += 10
	for i := 0; i < 9; i++ {
		b.Matrix[i] = int32(binary.BigEndian.Uint32(b.raw[offset+i : offset+i+4]))
	}
	offset += (9 * 4)
	b.Predefined = b.raw[offset : offset+24]
	offset += 24
	b.NextTrackID = binary.BigEndian.Uint32(b.raw[offset : offset+4])

	b.raw = nil
	return nil
}
