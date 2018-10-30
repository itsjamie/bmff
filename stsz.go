package bmff

import "encoding/binary"

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
