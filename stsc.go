package bmff

import (
	"encoding/binary"
)

type SampleToChunk struct {
	*fullbox
	Entries []SampleToChunkEntry
}

func (b *SampleToChunk) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	offset := 4
	b.Entries = make([]SampleToChunkEntry, 0, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		entry := SampleToChunkEntry{
			firstChunk:             binary.BigEndian.Uint32(b.raw[offset : offset+4]),
			samplesPerChunk:        binary.BigEndian.Uint32(b.raw[offset+4 : offset+8]),
			sampleDescriptionIndex: binary.BigEndian.Uint32(b.raw[offset+8 : offset+12]),
		}
		b.Entries = append(b.Entries, entry)
		offset += 12
	}

	return nil
}

type SampleToChunkEntry struct {
	firstChunk             uint32
	samplesPerChunk        uint32
	sampleDescriptionIndex uint32
}
