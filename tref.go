package bmff

import (
	"encoding/binary"
)

// TrackReference provides a reference from the containing track to another track in the presentation.
type TrackReference struct {
	*box
	TypeBoxes []*TrackReferenceType
}

func (b *TrackReference) parse() error {
	for subBox := range readBoxes(b.raw) {
		t := TrackReferenceType{box: subBox}
		for i := 0; i < len(t.raw); i += 4 {
			t.TrackIDs = append(t.TrackIDs, binary.BigEndian.Uint32(t.raw[i:i+4]))

		}
		b.TypeBoxes = append(b.TypeBoxes, &t)
	}

	b.raw = nil
	return nil
}

// TrackReferenceType is a child from the `tref` box. It contains information such as:
// - `hint` will reference links from the containing hint track to the media data it hints.
// - `cdsc` links a descriptive or metadata track to the content it describes
// See ISO 14936-12 8.3.3.3 for additional definitions.
type TrackReferenceType struct {
	*box
	TrackIDs []uint32
}
