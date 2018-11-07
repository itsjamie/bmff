package bmff

import (
	"fmt"
)

type TrakBox struct {
	*box
	Header    *TrackHeader
	Reference *TrackReference
	Media     *Media
	EditBox   *EditBox
	UserData  *UserData
}

func (b *TrakBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.boxtype {
		case "tkhd":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "tkhd":
			tkhd := &TrackHeader{fullbox: fb}
			if err := tkhd.parse(); err != nil {
				return err
			}
			b.Header = tkhd
		case "tref":
			tref := &TrackReference{box: subBox}
			if err := tref.parse(); err != nil {
				return err
			}
			b.Reference = tref
		case "mdia":
			mdia := &Media{box: subBox}
			if err := mdia.parse(); err != nil {
				return err
			}
			b.Media = mdia
		case "edts":
			edts := &EditBox{box: subBox}
			if err := edts.parse(); err != nil {
				return err
			}
			b.EditBox = edts
		case "udta":
			udta := &UserData{box: subBox}
			if err := udta.parse(); err != nil {
				return err
			}
			b.UserData = udta
		default:
			return fmt.Errorf("unknown 'trak' child: %s", subBox.Type())
		}

	}

	b.raw = nil
	return nil
}
