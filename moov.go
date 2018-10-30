package bmff

import (
	"fmt"
)

// The Movie Box is short-formed as the type 'moov'.
// The metadata for a presentation is stored in the single Movie Box which occurs at the top‐level of a file.
// Normally this box is close to the beginning or end of the file, though this is not required.
type Movie struct {
	*box
	Header   *MovieHeader
	Tracks   []*TrakBox
	Iods     *IodsBox
	Metadata *box
	UserData *UserData
	Unknown  []*box
}

func (b *Movie) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.boxtype {
		case "mvhd", "iods":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "mvhd":
			b.Header = &MovieHeader{fullbox: fb}
			if err := b.Header.parse(); err != nil {
				return err
			}
		case "iods":
			b.Iods = &IodsBox{fullbox: fb}
			if err := b.Iods.parse(); err != nil {
				return err
			}
		case "trak":
			trak := &TrakBox{box: subBox}
			if err := trak.parse(); err != nil {
				return err
			}
			b.Tracks = append(b.Tracks, trak)
		case "udta":
			udta := &UserData{box: subBox}
			if err := udta.parse(); err != nil {
				return err
			}
			b.UserData = udta
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown '%s' child: %s\n", b.boxtype, subBox.Type())

		}
	}

	b.raw = nil
	return nil
}
