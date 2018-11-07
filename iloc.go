package bmff

import (
	"encoding/binary"

	"github.com/itsjamie/bmff/internal/mp4bits"
)

type ItemLocation struct {
	*fullbox
	offsetSize     uint8
	lengthSize     uint8
	baseOffsetSize uint8
	indexSize      uint8
	itemCount      uint32
	Items          []ItemEntry
}

func (b *ItemLocation) parse() error {
	hb := mp4bits.HalvedByte(b.raw[0])
	b.offsetSize = hb.High()
	b.lengthSize = hb.Low()
	hb = mp4bits.HalvedByte(b.raw[1])
	b.baseOffsetSize = hb.High()
	if b.version == 1 || b.version == 2 {
		b.indexSize = hb.Low()
	}
	offset := 6
	if b.version == 2 {
		b.itemCount = binary.BigEndian.Uint32(b.raw[2:6])
	} else {
		b.itemCount = uint32(binary.BigEndian.Uint16(b.raw[2:4]))
		offset = 4
	}

	b.Items = make([]ItemEntry, 0, b.itemCount)
	for i := 0; uint32(i) < b.itemCount; i++ {
		item := ItemEntry{}
		if b.version == 2 {
			item.ID = binary.BigEndian.Uint32(b.raw[offset : offset+4])
			offset += 4
		} else {
			item.ID = uint32(binary.BigEndian.Uint16(b.raw[offset : offset+2]))
			offset += 2
		}

		if b.version == 1 || b.version == 2 {
			hb = mp4bits.HalvedByte(b.raw[offset+1]) // skip one byte, reserved
			item.constructionMethod = hb.Low()
			offset += 2
		}

		item.dataReferenceIndex = binary.BigEndian.Uint16(b.raw[offset : offset+2])
		offset += 2

		switch b.baseOffsetSize {
		case 4:
			item.baseOffset = uint64(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
			offset += 4
		case 8:
			item.baseOffset = uint64(binary.BigEndian.Uint64(b.raw[offset : offset+8]))
			offset += 8
		}

		extentCount := binary.BigEndian.Uint16(b.raw[offset : offset+2])
		offset += 2
		item.Extents = make([]Extent, 0, extentCount)
		for j := 0; uint16(j) < extentCount; j++ {
			extent := Extent{}
			if (b.version == 1 || b.version == 2) && b.indexSize > 0 {
				switch b.indexSize {
				case 4:
					extent.index = uint64(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
					offset += 4
				case 8:
					extent.index = uint64(binary.BigEndian.Uint64(b.raw[offset : offset+8]))
					offset += 8
				}
			}

			switch b.offsetSize {
			case 4:
				extent.offset = uint64(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
				offset += 4
			case 8:
				extent.offset = uint64(binary.BigEndian.Uint64(b.raw[offset : offset+8]))
				offset += 8
			}

			switch b.lengthSize {
			case 4:
				extent.length = uint64(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
				offset += 4
			case 8:
				extent.length = uint64(binary.BigEndian.Uint64(b.raw[offset : offset+8]))
				offset += 8
			}

			item.Extents = append(item.Extents, extent)
		}
		b.Items = append(b.Items, item)
	}
	b.raw = nil
	return nil
}

type ItemEntry struct {
	ID                 uint32
	constructionMethod uint8
	dataReferenceIndex uint16
	baseOffset         uint64
	Extents            []Extent
}

type Extent struct {
	index  uint64
	offset uint64
	length uint64
}
