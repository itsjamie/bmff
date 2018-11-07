package bmff

import (
	"encoding/binary"
	"fmt"
)

type EditBox struct {
	*box
	EditList *EditList
}

func (b *EditBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.Type() {
		case "elst":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.Type() {
		case "elst":
			elst := &EditList{fullbox: fb}
			if err := elst.parse(); err != nil {
				return err
			}
			b.EditList = elst
		default:
			return fmt.Errorf("unknown '%s' child: %s", b.boxtype, subBox.Type())
		}
	}

	b.raw = nil
	return nil
}

type EditList struct {
	*fullbox
	Versioned         EditLister
	MediaRate         int16
	MediaRateFraction int16
}

type EditLister interface {
	SegmentDuration() uint64
	MediaTime() int64
}

type EditListV0 struct {
	segmentDuration uint32
	mediaTime       int32
}

func (b *EditListV0) SegmentDuration() uint64 {
	return uint64(b.segmentDuration)
}

func (b *EditListV0) MediaTime() int64 {
	return int64(b.mediaTime)
}

type EditListV1 struct {
	segmentDuration uint64
	mediaTime       int64
}

func (b *EditListV1) SegmentDuration() uint64 {
	return b.segmentDuration
}

func (b *EditListV1) MediaTime() int64 {
	return b.mediaTime
}

func (b *EditList) parse() error {
	entryCount := binary.BigEndian.Uint32(b.raw[0:4])
	offset := 4
	for i := 0; uint32(i) < entryCount; i++ {
		if b.version == 1 {
			b.Versioned = &EditListV1{
				segmentDuration: binary.BigEndian.Uint64(b.raw[offset : offset+8]),
				mediaTime:       int64(binary.BigEndian.Uint64(b.raw[offset+8 : offset+16])),
			}

			offset += 16
		} else {
			b.Versioned = &EditListV0{
				segmentDuration: binary.BigEndian.Uint32(b.raw[offset : offset+4]),
				mediaTime:       int32(binary.BigEndian.Uint32(b.raw[offset+4 : offset+8])),
			}

			offset += 8
		}
		b.MediaRate = int16(binary.BigEndian.Uint16(b.raw[offset : offset+2]))
		b.MediaRateFraction = int16(binary.BigEndian.Uint16(b.raw[offset+2 : offset+4]))
		offset += 4
	}

	b.raw = nil
	return nil
}
