package bmff

import (
	"encoding/binary"
	"time"

	"github.com/itsjamie/bmff/internal/fixed"
)

const trackEnabled = 0x000001
const trackInMovie = 0x000002
const trackInPreview = 0x000004
const trackSizeIsAspectRatio = 0x000008

type TrackHeader struct {
	*fullbox
	creationTime     uint64
	modificationTime uint64
	TrackID          uint32
	Duration         uint64
	Layer            int16
	AlternateGroup   int16
	Volume           int16
	Matrix           [9]int32
	Width            fixed.Uint16_16
	Height           fixed.Uint16_16
}

func (b *TrackHeader) parse() error {
	var offset int
	if b.version == 0 {
		b.creationTime = uint64(binary.BigEndian.Uint32(b.raw[0:4]))
		b.modificationTime = uint64(binary.BigEndian.Uint32(b.raw[4:8]))
		b.TrackID = binary.BigEndian.Uint32(b.raw[8:12])
		// 12:16 reserved
		b.Duration = uint64(binary.BigEndian.Uint32(b.raw[16:20]))
		offset = 20
	} else if b.version == 1 {
		b.creationTime = uint64(binary.BigEndian.Uint64(b.raw[0:8]))
		b.modificationTime = uint64(binary.BigEndian.Uint64(b.raw[8:16]))
		b.TrackID = binary.BigEndian.Uint32(b.raw[16:20])
		// 20:24 reserved
		b.Duration = uint64(binary.BigEndian.Uint64(b.raw[24:28]))
		offset = 32
	}
	offset += 8 // reserved bytes
	b.Layer = int16(binary.BigEndian.Uint16(b.raw[offset : offset+2]))
	offset += 2
	b.AlternateGroup = int16(binary.BigEndian.Uint16(b.raw[offset : offset+2]))
	offset += 2
	b.Volume = int16(binary.BigEndian.Uint16(b.raw[offset : offset+2]))
	offset += 4 // previous bytes + 2 reserved

	for i := 0; i < 9; i++ {
		b.Matrix[i] = int32(binary.BigEndian.Uint32(b.raw[offset+i : offset+i+4]))
	}
	offset += 36
	b.Width = fixed.Uint16_16(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
	offset += 4
	b.Height = fixed.Uint16_16(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
	b.raw = nil
	return nil
}

// Enabled indicates that the track is present.
// A track which is not enabled should be treated as not present.
func (b *TrackHeader) Enabled() bool {
	return (b.flags & trackEnabled) > 0
}

// InMovie indicates that the track is used in the presentation
func (b *TrackHeader) InMovie() bool {
	return (b.flags & trackInMovie) > 0
}

// InPreview indiciates the track is used when previewing the presentation
func (b *TrackHeader) InPreview() bool {
	return (b.flags & trackInPreview) > 0
}

// SizeIsAspectRatio indicates that the width and height values are not expressed in pixel units.
// The values have the same unit but the unit is not specified, the values are only an indication of the desired aspect ratio.
// If the aspect ratios of this track and other related tracks are not identical then the respective positioning is undefined, possibly defined by external contexts.
func (b *TrackHeader) SizeIsAspectRatio() bool {
	return (b.flags & trackSizeIsAspectRatio) > 0
}

// CreationTime is when this track was created
func (b *TrackHeader) CreationTime() time.Time {
	epoch := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	return epoch.Add(time.Duration(b.creationTime) * time.Second)
}

// ModificationTime is when this track was most recently edited
func (b *TrackHeader) ModificationTime() time.Time {
	epoch := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	return epoch.Add(time.Duration(b.modificationTime) * time.Second)
}
