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
		switch t := subBox.Type(); t {
		case "mdhd":
			mdhd := MediaHeader{box: subBox}
			if err := mdhd.parse(); err != nil {
				return err
			}
			b.Header = &mdhd
		case "hdlr":
			hdlr := Handler{box: subBox}
			if err := hdlr.parse(); err != nil {
				return err
			}
			b.Handler = &hdlr
		case "minf":
			minf := MediaInformation{box: subBox}
			if err := minf.parse(); err != nil {
				return err
			}
			b.Information = &minf
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Errorf("unknown mdia subtype: %s", t)
		}
	}

	// b.raw = nil
	return nil
}
