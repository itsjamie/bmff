package bmff

import (
	"github.com/itsjamie/bmff/internal/mp4bits"
)

type SampleDependenciesTable struct {
	*fullbox
	Entries []SampleDependencyEntry
}

func (s *SampleDependenciesTable) parse() error {
	s.Entries = make([]SampleDependencyEntry, 0, len(s.raw))

	for _, b := range s.raw {
		b := mp4bits.Int2(b)
		entry := SampleDependencyEntry{
			IsLeading:           b.First(),
			SampleDependsOn:     b.Second(),
			SampleIsDependedOn:  b.Third(),
			SampleHasRedundancy: b.Fourth(),
		}
		s.Entries = append(s.Entries, entry)
	}
	s.raw = nil
	return nil
}

type SampleDependencyEntry struct {
	IsLeading           uint8
	SampleDependsOn     uint8
	SampleIsDependedOn  uint8
	SampleHasRedundancy uint8
}
