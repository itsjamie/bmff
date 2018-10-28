package bmff

import (
	"fmt"
)

type SampleTable struct {
	*box
	Stsd    *SampleDescription
	Stts    *TimeToSample
	Stsc    *SampleToChunk
	Stsz    *SampleSize
	Stco    *ChunkOffset
	Co64    *ChunkLargeOffset
	Unknown []*box
}

func (b *SampleTable) parse() error {
	for subBox := range readBoxes(b.raw) {
		var fb *fullbox
		switch subBox.boxtype {
		case
			"stsd", // Sample Description
			"stts", // Time to Sample
			"stsc", // Sample to Chunk
			"stsz", // Sample Size
			"stco", // Chunk Offset 32 bit
			"co64": // Chunk Offset 64 bit
			fb = &fullbox{box: subBox}
			if err := fb.decode(); err != nil {
				return err
			}
		}

		switch subBox.boxtype {
		case "stsd":
			stsd := &SampleDescription{fullbox: fb}
			if err := stsd.parse(); err != nil {
				return err
			}
			b.Stsd = stsd
		case "stts":
			stts := &TimeToSample{fullbox: fb}
			if err := stts.parse(); err != nil {
				return err
			}
			b.Stts = stts
		case "stsc":
			stsc := &SampleToChunk{fullbox: fb}
			if err := stsc.parse(); err != nil {
				return err
			}
			b.Stsc = stsc
		case "stsz":
			stsz := &SampleSize{fullbox: fb}
			if err := stsz.parse(); err != nil {
				return err
			}
			b.Stsz = stsz
		case "stco":
			stco := &ChunkOffset{fullbox: fb}
			if err := stco.parse(); err != nil {
				return err
			}
			b.Stco = stco
		case "co64":
			co64 := &ChunkLargeOffset{fullbox: fb}
			if err := co64.parse(); err != nil {
				return err
			}
			b.Co64 = co64
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown '%s' child: %s\n", b.boxtype, subBox.Type())
		}
	}

	return nil
}
