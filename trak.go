package bmff

import (
	"fmt"
)

type TrakBox struct {
	*box
	Header    *TrackHeader
	Reference *TrefBox
	Media     *Media
}

func (b *TrakBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		if subBox == nil {
			break
		}

		switch t := subBox.boxtype; t {
		case "tkhd":
			header := &TrackHeader{box: subBox}
			if err := header.parse(); err != nil {
				return err
			}
			b.Header = header
		case "mdia":
			mdia := &Media{box: subBox}
			if err := mdia.parse(); err != nil {
				return err
			}
			b.Media = mdia
		case "tref":
			tref := &TrefBox{box: subBox}
			if err := tref.parse(); err != nil {
				return err
			}
			b.Reference = tref
		default:
			return fmt.Errorf("unknown box in trak: %s", t)
		}
	}

	return nil
}
