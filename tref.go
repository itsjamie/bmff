package bmff

import "encoding/binary"

type TrefBox struct {
	*box
	TypeBoxes []*TrefTypeBox
}

func (b *TrefBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		t := TrefTypeBox{box: subBox}
		for i := 0; i < len(t.raw); i += 4 {
			t.TrackIDs = append(t.TrackIDs, binary.BigEndian.Uint32(t.raw[i:i+4]))

		}
		b.TypeBoxes = append(b.TypeBoxes, &t)
	}

	return nil
}

type TrefTypeBox struct {
	*box
	TrackIDs []uint32
}
