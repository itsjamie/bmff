package bmff

import "fmt"

type Metadata struct {
	*fullbox
	Handler      *Handler
	ItemLocation *ItemLocation
	Unknown      []*box
}

func (b *Metadata) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.Type() {
		case "hdlr", "iloc":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.Type() {
		case "hdlr":
			hdlr := &Handler{fullbox: fb}
			if err := hdlr.parse(); err != nil {
				return err
			}
			b.Handler = hdlr
		case "iloc":
			iloc := &ItemLocation{fullbox: fb}
			if err := iloc.parse(); err != nil {
				return err
			}
			b.ItemLocation = iloc
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown '%s' child: %s\n", b.boxtype, subBox.Type())
		}
	}

	b.raw = nil
	return nil
}
