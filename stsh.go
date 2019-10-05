package bmff

import (
	"encoding/binary"
)

type ShadowSyncSample struct {
	*fullbox
	Entries []ShadowSyncEntry
}

type ShadowSyncEntry struct {
	ShadowedSampleNumber uint32
	SyncSampleNumber     uint32
}

func (b *ShadowSyncSample) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	offset := 4
	b.Entries = make([]ShadowSyncEntry, 0, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		entry := ShadowSyncEntry{
			ShadowedSampleNumber: binary.BigEndian.Uint32(b.raw[offset : offset+4]),
			SyncSampleNumber:     binary.BigEndian.Uint32(b.raw[offset+4 : offset+8]),
		}
		offset += 8
		b.Entries = append(b.Entries, entry)
	}
	return nil
}
