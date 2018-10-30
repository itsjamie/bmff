package bmff

import (
	"encoding/binary"
)

type ChunkOffset struct {
	*fullbox
	Entries []uint32
}

func (b *ChunkOffset) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	b.Entries = make([]uint32, 0, entryCount)
	offset := 4
	for i := 0; uint32(i) < entryCount; i++ {
		b.Entries = append(b.Entries, binary.BigEndian.Uint32(b.raw[offset:offset+4]))
		offset += 4
	}
	return nil
}

type ChunkLargeOffset struct {
	*fullbox
	Entries []uint64
}

func (b *ChunkLargeOffset) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	b.Entries = make([]uint64, 0, entryCount)
	offset := 4
	for i := 0; uint32(i) < entryCount; i++ {
		b.Entries = append(b.Entries, binary.BigEndian.Uint64(b.raw[offset:offset+8]))
		offset += 8
	}

	b.raw = nil
	return nil
}
