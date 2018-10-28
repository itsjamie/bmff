package bmff

import (
	"encoding/binary"
)

type TimeToSample struct {
	*fullbox
	Entries []TimeToSampleEntry
}

func (b *TimeToSample) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	b.Entries = make([]TimeToSampleEntry, 0, entryCount)

	offset := 4

	for i := 0; uint32(i) < entryCount; i++ {
		entry := TimeToSampleEntry{
			sampleCount: binary.BigEndian.Uint32(b.raw[offset : offset+4]),
			sampleDelta: binary.BigEndian.Uint32(b.raw[offset+4 : offset+8]),
		}
		offset += 8
		b.Entries = append(b.Entries, entry)
	}

	b.raw = nil
	return nil
}

type TimeToSampleEntry struct {
	sampleCount uint32
	sampleDelta uint32
}
