package bmff

import (
	"encoding/binary"
	"errors"
)

type SampleDescription struct {
	*fullbox
	entryCount uint32
	Entries    []*SampleEntry
}

func (b *SampleDescription) parse() error {
	b.entryCount = binary.BigEndian.Uint32(b.raw[0:4])
	b.Entries = make([]*SampleEntry, 0, b.entryCount)
	b.raw = b.raw[4:]
	for subBox := range readBoxes(b.raw) {
		entry := &SampleEntry{box: subBox}
		if err := entry.parse(); err != nil {
			return err
		}

		b.Entries = append(b.Entries, entry)
	}

	if uint32(len(b.Entries)) > b.entryCount {
		return errors.New("SampleEntries > entryCount")
	}

	b.raw = nil
	return nil
}

type SampleEntry struct {
	*box
	data_reference_index uint16
	sample               []byte
}

func (b *SampleEntry) parse() error {
	// first six bytes are reserved

	b.data_reference_index = binary.BigEndian.Uint16(b.raw[6:8])
	b.sample = b.raw[8:]
	b.raw = nil

	return nil
}
