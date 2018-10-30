package bmff

import (
	"encoding/binary"
)

type DataReference struct {
	*fullbox
	Entries []DataEntry
}

func (b *DataReference) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	b.raw = b.raw[4:]

	b.Entries = make([]DataEntry, 0, entryCount)

	for subBox := range readBoxes(b.raw) {

		var fb *fullbox
		switch subBox.boxtype {
		case "url ", "urn ":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "url ":
			url := &DataEntryURL{fullbox: fb}
			if err := url.parse(); err != nil {
				return err
			}
			b.Entries = append(b.Entries, url)
		case "urn ":
			urn := &DataEntryURN{fullbox: fb}
			if err := urn.parse(); err != nil {
				return err
			}
			b.Entries = append(b.Entries, urn)
		}
	}

	b.raw = nil
	return nil
}

const inFile = 0x01

type DataEntry interface {
	Name() string
	Location() string
	InFile() bool
}

type DataEntryURL struct {
	*fullbox
	location string
}

func (b *DataEntryURL) parse() error {
	if !b.InFile() {
		b.location = string(b.raw[:])
	}

	return nil
}

func (b *DataEntryURL) Name() string {
	return ""
}

func (b *DataEntryURL) Location() string {
	return b.location
}

func (b *DataEntryURL) InFile() bool {
	return (b.flags & inFile) > 0
}

type DataEntryURN struct {
	*fullbox
	name     string
	location string
}

func (b *DataEntryURN) parse() error {
	firstTerminator := clen(b.raw[:])
	b.name = string(b.raw[0:firstTerminator])
	b.location = string(b.raw[firstTerminator:])

	return nil
}

func (b *DataEntryURN) Name() string {
	return b.name
}

func (b *DataEntryURN) Location() string {
	return b.location
}

func (b *DataEntryURN) InFile() bool {
	return (b.flags & inFile) > 0
}
