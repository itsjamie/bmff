package bmff

import "fmt"

type DataInformation struct {
	*box
}

func (b *DataInformation) parse() error {
	return fmt.Errorf("not implemented")
}
