package bmff

import "fmt"

type DataInformation struct {
	*box
	Reference *DataReference
	Unknown   []*box
}

func (b *DataInformation) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.boxtype {
		case "dref":
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "dref":
			ref := &DataReference{fullbox: fb}
			if err := ref.parse(); err != nil {
				return err
			}
			b.Reference = ref
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown '%s' child: %s\n", b.boxtype, subBox.Type())
		}
	}

	b.raw = nil
	return nil
}
