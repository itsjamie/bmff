package bmff

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type File struct {
	Type      *FileType
	Movie     *Movie
	MediaData []*box
	Metadata  *Metadata
	Free      []*box
	Unknown   []*box

	freeSkipBeforeMdat bool
}

func Parse(src io.Reader) (*File, error) {
	f := &File{}
	r := bufio.NewReader(src)
	var parseFreeOrSkip bool

readloop:
	for {
		b, err := newBox(r)
		if err != nil {
			switch errors.Cause(err) {
			case io.EOF:
				if b == nil {
					break readloop
				}
			default:
				return nil, err
			}
		}

		switch t := b.Type(); t {
		case "ftyp":
			fb := &FileType{box: b}
			if err := fb.parse(); err != nil {
				return nil, err
			}
			f.Type = fb
		case "moov":
			mb := &Movie{box: b}
			if err := mb.parse(); err != nil {
				return nil, err
			}
			f.Movie = mb
		case "mdat":
			if parseFreeOrSkip {
				f.freeSkipBeforeMdat = true
			}

			f.MediaData = append(f.MediaData, b)
		case "meta":
			fb := &fullbox{box: b}
			if err := fb.decode(); err != nil {
				return nil, err
			}

			meta := &Metadata{fullbox: fb}
			if err := meta.parse(); err != nil {
				return nil, err
			}
			f.Metadata = meta
		case "skip", "free":
			f.Free = append(f.Free, b)
			parseFreeOrSkip = true
		default:
			f.Unknown = append(f.Unknown, b)
			fmt.Printf("unknown top-level box: %s\n", b.Type())
		}
	}

	return f, nil
}

// returns the length until first null terminated byte for c-strings
func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
