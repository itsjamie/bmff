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
		switch subBox.boxtype {
		case "nmhd":
			nmhd := &NullMediaHeader{
				box: subBox,
			}
			if err := nmhd.parse(); err != nil {
				return err
			}
			if b.MediaHeader != nil {
				return fmt.Errorf("media header already populated for track: %v", b.MediaHeader)
			}
			b.MediaHeader = nmhd
		case "dinf":
			dinf := &DataInformation{
				box: subBox,
			}
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
			fmt.Printf("unknown minf subtype: %s\n", subBox.boxtype)

		}
	}

	b.raw = nil
	return nil
}

type NullMediaHeader struct {
	*box
	version uint8
	flags   [3]byte
}

func (b *NullMediaHeader) parse() error {
	b.version = b.raw[0]
	b.flags = [3]byte{b.raw[1], b.raw[2], b.raw[3]}
	return nil
}

type VideoMediaHeader struct {
}
