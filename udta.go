package bmff

import (
	"encoding/binary"

	"github.com/davecgh/go-spew/spew"
)

type UserData struct {
	*box
	Copyright *Copyright
	Boxes     []*box
}

func (b *UserData) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.boxtype {
		case "cprt":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "cprt":
			cprt := &Copyright{fullbox: fb}
			if err := cprt.parse(); err != nil {
				return err
			}
			b.Copyright = cprt
		default:
			b.Boxes = append(b.Boxes, subBox)
		}

	}
	spew.Dump(b)
	return nil
}

type Copyright struct {
	*fullbox
	LanguageCode string
	Data         string
}

func (b *Copyright) parse() error {
	lang := binary.BigEndian.Uint16(b.raw[0:2]) // first bit is padding
	b.LanguageCode = string([]byte{
		uint8(lang&0x7C00>>10) + 0x60,
		uint8(lang&0x03E0>>5) + 0x60,
		uint8(lang&0x001F) + 0x60,
	})
	b.Data = string(b.raw[2:])
	spew.Dump(b)
	return nil
}
