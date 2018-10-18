package bmff

import (
	"encoding/binary"
	"fmt"
)

type SampleTable struct {
	*box
	Unknown []*box
}

func (b *SampleTable) parse() error {
	for subBox := range readBoxes(b.raw) {
		switch subBox.boxtype {
		case "stsd":
			stsd := &SampleDescription{box: subBox}
			if err := stsd.parse(); err != nil {
				return err
			}
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown stbl subtype: %s\n", subBox.boxtype)
		}
	}

	return nil
}

type SampleDescription struct {
	*box
	version    uint8
	flags      [3]byte
	entryCount uint32
	Entries    []*SampleEntry
}

func (b *SampleDescription) parse() error {
	b.version = b.raw[0]
	b.flags = [3]byte{b.raw[1], b.raw[2], b.raw[3]}
	b.entryCount = binary.BigEndian.Uint32(b.raw[4:8])

	b.Entries = make([]*SampleEntry, 0, b.entryCount)

	for subBox := range readBoxes(b.raw[8:]) {

		entry := &SampleEntry{box: subBox}
		if err := entry.parse(); err != nil {
			return err
		}
		b.Entries = append(b.Entries, entry)

		if uint32(len(b.Entries)) >= b.entryCount {
			break
		}
	}

	return nil
}

type SampleEntry struct {
	*box
	reserved             [6]uint8
	data_reference_index uint16
}

func (b *SampleEntry) parse() error {
	// first six bytes are reserved
	b.data_reference_index = binary.BigEndian.Uint16(b.raw[6:8])
	b.raw = b.raw[8:]

	return nil
}
