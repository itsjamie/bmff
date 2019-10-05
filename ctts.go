package bmff

import (
	"encoding/binary"
	"log"
)

// CompositionOffset provides the offset between decoding time and composition time.
type CompositionOffset struct {
	*fullbox
	Entries []CompositionOffsetEntry
}

// CompositionOffsetEntry represents a single entry
type CompositionOffsetEntry interface {
	Samples() uint32
	Offset() int32
}

// CompositionOffsetEntryV0 represents a V0 entry, offset is unsigned
type CompositionOffsetEntryV0 struct {
	sampleCount  uint32
	sampleOffset uint32
}

// Samples returns the number of samples
func (e CompositionOffsetEntryV0) Samples() uint32 {
	return e.sampleCount
}

// Offset returns the offset
func (e CompositionOffsetEntryV0) Offset() int32 {
	return int32(e.sampleOffset)
}

// CompositionOffsetEntryV1 represents a V1 entry, offset is signed
type CompositionOffsetEntryV1 struct {
	sampleCount  uint32
	sampleOffset int32
}

// Samples returns the number of samples
func (e CompositionOffsetEntryV1) Samples() uint32 {
	return e.sampleCount
}

// Offset returns the offset
func (e CompositionOffsetEntryV1) Offset() int32 {
	return e.sampleOffset
}

func (b *CompositionOffset) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	offset := 4

	b.Entries = make([]CompositionOffsetEntry, 0, entryCount)
	switch b.version {
	case 0:
		for i := 0; uint32(i) < entryCount; i++ {
			entry := CompositionOffsetEntryV0{
				sampleCount:  binary.BigEndian.Uint32(b.raw[offset : offset+4]),
				sampleOffset: binary.BigEndian.Uint32(b.raw[offset+4 : offset+8]),
			}
			b.Entries = append(b.Entries, entry)
			offset += 8
		}
	case 1:
		for i := 0; uint32(i) < entryCount; i++ {
			entry := CompositionOffsetEntryV1{
				sampleCount:  binary.BigEndian.Uint32(b.raw[offset : offset+4]),
				sampleOffset: int32(binary.BigEndian.Uint32(b.raw[offset+4 : offset+8])),
			}
			b.Entries = append(b.Entries, entry)
		}
	default:
		log.Fatal("unknown ctts version: ", b.version)
	}

	b.raw = nil
	return nil
}
