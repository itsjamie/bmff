package bmff

import (
	"encoding/binary"
	"time"
)

type MediaHeader struct {
	*fullbox
	creationTime     uint64
	modificationTime uint64
	TimeScale        uint32
	Duration         uint64
	LanguageCode     string
	Predefined       uint16
}

func (b *MediaHeader) parse() error {
	var offset int
	raw := b.Raw()
	if b.version == 0 {
		b.creationTime = uint64(binary.BigEndian.Uint32(raw[0:4]))
		b.modificationTime = uint64(binary.BigEndian.Uint32(raw[4:8]))
		b.TimeScale = binary.BigEndian.Uint32(raw[8:12])
		b.Duration = uint64(binary.BigEndian.Uint32(raw[12:16]))
		offset = 16
	} else if b.version == 1 {
		b.creationTime = binary.BigEndian.Uint64(raw[0:8])
		b.modificationTime = binary.BigEndian.Uint64(raw[8:16])
		b.TimeScale = binary.BigEndian.Uint32(raw[16:20])
		b.Duration = binary.BigEndian.Uint64(raw[20:28])
		offset = 28
	}

	lang := binary.BigEndian.Uint16(raw[offset : offset+2]) // first bit is padding
	b.LanguageCode = string([]byte{
		uint8(lang&0x7C00>>10) + 0x60,
		uint8(lang&0x03E0>>5) + 0x60,
		uint8(lang&0x001F) + 0x60,
	})
	b.Predefined = binary.BigEndian.Uint16(raw[offset+2 : offset+4])

	b.raw = nil
	return nil
}

// CreationTime is when this track was created
func (b *MediaHeader) CreationTime() time.Time {
	epoch := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	return epoch.Add(time.Duration(b.creationTime) * time.Second)
}

// ModificationTime is when this track was most recently edited
func (b *MediaHeader) ModificationTime() time.Time {
	epoch := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	return epoch.Add(time.Duration(b.modificationTime) * time.Second)
}
