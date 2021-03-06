package bmff

import (
	"fmt"
)

type SampleTable struct {
	*box
	Ctts    *CompositionOffset
	Stsd    *SampleDescription
	Sdtp    *SampleDependenciesTable
	Stts    *TimeToSample
	Stsc    *SampleToChunk
	Stss    *SyncSample
	Stsz    *SampleSize
	Stz2    *CompactSampleSize
	Stsh    *ShadowSyncSample
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
			"sdtp", // Sample Dependencies Type
			"stss", // Sync Sample
			"stts", // Time to Sample
			"stsc", // Sample to Chunk
			"stsz", // Sample Size
			"stz2", // Compact Sample Size
			"stsh", // ShadowSync
			"stco", // Chunk Offset 32 bit
			"co64", // Chunk Offset 64 bit
			"ctts": // Composition Offset
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
		case "sdtp":
			sdtp := &SampleDependenciesTable{fullbox: fb}
			if err := sdtp.parse(); err != nil {
				return err
			}
			b.Sdtp = sdtp
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
		case "stss":
			stss := &SyncSample{fullbox: fb}
			if err := stss.parse(); err != nil {
				return err
			}
			b.Stss = stss
		case "stsz":
			stsz := &SampleSize{fullbox: fb}
			if err := stsz.parse(); err != nil {
				return err
			}
			b.Stsz = stsz
		case "stz2":
			stz2 := &CompactSampleSize{fullbox: fb}
			if err := stz2.parse(); err != nil {
				return err
			}
			b.Stz2 = stz2
		case "stsh":
			stsh := &ShadowSyncSample{fullbox: fb}
			if err := stsh.parse(); err != nil {
				return err
			}
			b.Stsh = stsh
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
		case "ctts":
			ctts := &CompositionOffset{fullbox: fb}
			if err := ctts.parse(); err != nil {
				return err
			}
			b.Ctts = ctts
		default:
			b.Unknown = append(b.Unknown, subBox)
			fmt.Printf("unknown '%s' child: %s\n", b.boxtype, subBox.Type())
		}
	}

	b.raw = nil
	return nil
}
