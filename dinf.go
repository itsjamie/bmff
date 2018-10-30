package bmff

type DataInformation struct {
	*box
	Reference *DataReference
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
		}
	}

	b.raw = nil
	return nil
}
