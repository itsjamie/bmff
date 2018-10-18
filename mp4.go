package bmff

import (
	"bufio"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

type File struct {
	Type    *FileType
	Movie   *Movie
	Unknown []*box
}

func Parse(src io.Reader) (*File, error) {
	f := &File{}
	r := bufio.NewReader(src)

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
		default:
			f.Unknown = append(f.Unknown, b)
		}
	}

	spew.Dump(f)

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
