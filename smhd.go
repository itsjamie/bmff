package bmff

import (
	"github.com/itsjamie/bmff/internal/fixed"
)

type SoundMediaHeader struct {
	*fullbox
	balance fixed.Uint8_8
}

func (b *SoundMediaHeader) parse() error {
	return b.balance.UnmarshalBinary(b.raw[0:2])
}
