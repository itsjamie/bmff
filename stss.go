package bmff

import (
	"encoding/binary"
)

type SyncSample struct {
	*fullbox
	Entries []uint32
}

func (b *SyncSample) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	offset := 4
	b.Entries = make([]uint32, 0, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		b.Entries = append(b.Entries, binary.BigEndian.Uint32(b.raw[offset:offset+4]))
		offset += 4
	}
	return nil
}
