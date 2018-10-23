package bmff

import (
	"fmt"
)

type MediaInformation struct {
	*box
	DataInformation *DataInformation
	MediaHeader     *NullMediaHeader
	SampleTable     *SampleTable
	Unknown         []*box
}

func (b *MediaInformation) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.boxtype {
		case "nmhd":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "nmhd":
			nmhd := &NullMediaHeader{fullbox: fb}
			if err := nmhd.parse(); err != nil {
				return err
			}
			if b.MediaHeader != nil {
				return fmt.Errorf("media header already populated for track: %v", b.MediaHeader)
			}
			b.MediaHeader = nmhd
		case "dinf":
			dinf := &DataInformation{box: subBox}
			if err := dinf.parse(); err != nil {
				return err
			}
			b.DataInformation = dinf
		case "stbl":
			stbl := &SampleTable{
				box: subBox,
			}
			if err := stbl.parse(); err != nil {
				return err
			}
			b.SampleTable = stbl
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown '%s' child: %s\n", b.boxtype, subBox.Type())
		}
	}

	b.raw = nil
	return nil
}

type NullMediaHeader struct {
	*fullbox
}

func (b *NullMediaHeader) parse() error {
	return nil
}

type VideoMediaHeader struct {
}
