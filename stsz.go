package bmff

import (
	"encoding/binary"
	"errors"

	"github.com/itsjamie/bmff/internal/mp4bits"
)

type SampleSize struct {
	*fullbox
	DefaultSize uint32
	SampleCount uint32
	Entries     []uint32
}

func (b *SampleSize) parse() error {
	b.DefaultSize = binary.BigEndian.Uint32(b.raw[0:4])
	b.SampleCount = binary.BigEndian.Uint32(b.raw[4:8])

	if b.DefaultSize != 0 {
		return nil
	}

	b.Entries = make([]uint32, 0, b.SampleCount)
	offset := 8
	for i := 0; uint32(i) < b.SampleCount; i++ {
		b.Entries = append(b.Entries, binary.BigEndian.Uint32(b.raw[offset:offset+4]))
		offset += 4
	}

	b.raw = nil
	return nil
}

type CompactSampleSize struct {
	*fullbox
	FieldSize uint8
	Entries   []uint16
}

func (b *CompactSampleSize) parse() error {
	b.FieldSize = uint8(b.raw[3])
	sampleCount := binary.BigEndian.Uint32(b.raw[4:8])
	offset := 8

	for i := 0; uint32(i) < sampleCount; i++ {
		switch b.FieldSize {
		case 4:
			hb := mp4bits.HalvedByte(b.raw[offset])

			b.Entries = append(b.Entries, uint16(hb.High()))
			b.Entries = append(b.Entries, uint16(hb.Low()))
			offset++
		case 8:
			b.Entries = append(b.Entries, uint16(b.raw[offset]))
			offset++
		case 16:
			b.Entries = append(b.Entries, binary.BigEndian.Uint16(b.raw[offset:offset+2]))
			offset += 2
		}
	}

	if b.FieldSize == 4 && sampleCount%2 != 0 {
		if b.Entries[len(b.Entries)-1] != 0 {
			return errors.New("last entry should have been padded with zeroes")
		}
		b.Entries = b.Entries[0 : len(b.Entries)-1]
	}

	return nil
}
