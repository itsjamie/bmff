package bmff

import (
	"fmt"
)

// The Movie Box is short-formed as the type 'moov'.
// The metadata for a presentation is stored in the single Movie Box which occurs at the top‐level of a file.
// Normally this box is close to the beginning or end of the file, though this is not required.
type Movie struct {
	*box
	Header *MovieHeader
	Tracks []*TrakBox
	Iods   *IodsBox
	// Meta *MetaBox
	Unknown []*box
}

func (b *Movie) parse() error {
	for subBox := range readBoxes(b.raw) {
		switch subBox.boxtype {
		case "mvhd":
			b.Header = &MovieHeader{box: subBox}
			if err := b.Header.parse(); err != nil {
				return err
			}
		case "iods":
			b.Iods = &IodsBox{box: subBox}
			if err := b.Iods.parse(); err != nil {
				return err
			}
		case "trak":
			trak := &TrakBox{box: subBox}
			if err := trak.parse(); err != nil {
				return err
			}

			b.Tracks = append(b.Tracks, trak)
		default:
			fmt.Printf("Unknown Moov SubType: %s\n", subBox.Type())
			b.Unknown = append(b.Unknown, subBox)
		}
	}

	b.box.raw = nil
	return nil
}
