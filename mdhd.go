package bmff

import "encoding/binary"

type MediaHeader struct {
	*fullbox
	CreationTime     uint64
	ModificationTime uint64
	TimeScale        uint32
	Duration         uint64
	LanguageCode     string
	Predefined       uint16
}

func (b *MediaHeader) parse() error {
	var offset int
	raw := b.Raw()
	if b.version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(raw[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(raw[8:12]))
		b.TimeScale = binary.BigEndian.Uint32(raw[12:16])
		b.Duration = uint64(binary.BigEndian.Uint32(raw[16:20]))
		offset = 20
	} else if b.version == 1 {
		b.CreationTime = binary.BigEndian.Uint64(raw[4:12])
		b.ModificationTime = binary.BigEndian.Uint64(raw[12:20])
		b.TimeScale = binary.BigEndian.Uint32(raw[20:24])
		b.Duration = binary.BigEndian.Uint64(raw[24:32])
		offset = 32
	}

	lang := binary.BigEndian.Uint16(raw[offset : offset+2]) // first bit is padding
	b.LanguageCode = string([]byte{
		uint8(lang&0x7C00>>10) + 0x60,
		uint8(lang&0x03E0>>5) + 0x60,
		uint8(lang&0x001F) + 0x60,
	})
	b.Predefined = binary.BigEndian.Uint16(raw[offset+2 : offset+4])

	return nil
}
