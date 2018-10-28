package bmff

import (
	"encoding/binary"
)

// Handler is short-formed as 'hdlr'.
// This box within a Media Box declares media type of the track, and thus the process by which the media‐data in the track is presented.
// For example, a format for which the decoder delivers video would be stored in a video track, identified by being handled by a video handler.
// The documentation of the storage of a media format identifies the media type which that format uses.
// This box when present within a Meta Box, declares the structure or format of the 'meta' box contents.
// There is a general handler for metadata streams of any type; the specific format is identified by the sample entry, as for video or audio, for example.
type Handler struct {
	*fullbox
	Predefined  uint32
	HandlerType string
	Name        string
}

func (b *Handler) parse() error {
	b.Predefined = binary.BigEndian.Uint32(b.raw[0:4])
	b.HandlerType = string(b.raw[4:8])
	// skip 12 reserved bytes
	nameLen := clen(b.raw[20:])
	b.Name = string(b.raw[20 : 20+nameLen])

	b.raw = nil
	return nil
}
