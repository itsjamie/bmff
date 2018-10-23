package bmff

import "fmt"

type Media struct {
	*box
	Header      *MediaHeader
	Handler     *Handler
	Information *MediaInformation
	Unknown     []*box
}

func (b *Media) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.boxtype {
		case "mdhd", "hdlr":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "mdhd":
			mdhd := &MediaHeader{fullbox: fb}
			if err := mdhd.parse(); err != nil {
				return err
			}
			b.Header = mdhd
		case "hdlr":
			hdlr := &Handler{fullbox: fb}
			if err := hdlr.parse(); err != nil {
				return err
			}
			b.Handler = hdlr
		case "minf":
			minf := &MediaInformation{box: subBox}
			if err := minf.parse(); err != nil {
				return err
			}
			b.Information = minf
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown '%s' child: %s\n", b.boxtype, subBox.Type())
		}
	}

	// b.raw = nil
	return nil
}
