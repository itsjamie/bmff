package bmff

import (
	"encoding/binary"
	"log"
)

type CompositionOffset struct {
	*fullbox
	EntriesV0 []CompositionOffsetEntryV0
	EntriesV1 []CompositionOffsetEntryV1
}

type CompositionOffsetEntryV0 struct {
	sampleCount  uint32
	sampleOffset uint32
}

type CompositionOffsetEntryV1 struct {
	sampleCount  uint32
	sampleOffset int32
}

func (b *CompositionOffset) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	offset := 4

	switch b.version {
	case 0:
		b.EntriesV0 = make([]CompositionOffsetEntryV0, 0, entryCount)
		for i := 0; uint32(i) < entryCount; i++ {
			entry := CompositionOffsetEntryV0{
				sampleCount:  binary.BigEndian.Uint32(b.raw[offset : offset+4]),
				sampleOffset: binary.BigEndian.Uint32(b.raw[offset+4 : offset+8]),
			}
			b.EntriesV0 = append(b.EntriesV0, entry)
			offset += 8
		}
	case 1:
		b.EntriesV1 = make([]CompositionOffsetEntryV1, 0, entryCount)
		for i := 0; uint32(i) < entryCount; i++ {
			entry := CompositionOffsetEntryV1{
				sampleCount:  binary.BigEndian.Uint32(b.raw[offset : offset+4]),
				sampleOffset: int32(binary.BigEndian.Uint32(b.raw[offset+4 : offset+8])),
			}
			b.EntriesV1 = append(b.EntriesV1, entry)
		}
	default:
		log.Fatal("unknown ctts version: ", b.version)
	}

	return nil
}
