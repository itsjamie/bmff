package bmff

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

const boxHeaderSize = 8

func newBox(src io.Reader) (*box, error) {
	b := &box{}
	err := b.decode(src)

	return b, err
}

type box struct {
	boxtype   string
	size      uint32
	largesize uint64
	raw       []byte
}

type fullbox struct {
	box
	version uint8
	flags   [3]byte
	raw     []byte
}

func (b *box) decode(src io.Reader) error {
	header := make([]byte, 8)
	if _, err := io.ReadFull(src, header); err != nil {
		return errors.Wrapf(err, "bad read of box header: %x", header)
	}

	b.size = binary.BigEndian.Uint32(header[0:4])
	b.boxtype = string(header[4:8])

	if b.size == 1 {
		if err := binary.Read(src, binary.BigEndian, &b.largesize); err != nil {
			return errors.Wrap(err, "failure to read largesize field")
		}
	}

	if b.size == 1 {
		b.raw = make([]byte, b.largesize-(boxHeaderSize-8))
		if _, err := io.ReadFull(src, b.raw); err != nil {
			return errors.Wrapf(err, "bad read for box(%s) body", b.boxtype)
		}
	} else if b.size == 0 {
		data, err := ioutil.ReadAll(src)
		if err != nil {
			return errors.Wrap(err, "failure reading box to EOF")
		}
		b.raw = data
	} else {
		size := b.size - boxHeaderSize
		if b.largesize != 0 {
			size -= 8
		}

		b.raw = make([]byte, size)
		if _, err := io.ReadFull(src, b.raw); err != nil {
			return errors.Wrapf(err, "bad read for box(%s) body", b.boxtype)
		}
	}

	spew.Dump(b)

	return nil
}

func (b *fullbox) decode(src *bufio.Reader) error {
	b.box.decode(src)

	if len(b.raw) < 4 {
		return errors.Errorf("failed decode of fullbox type %s, missing data", b.box.boxtype)
	}

	b.version = uint8(b.raw[0])
	b.flags = [3]byte{b.raw[1], b.raw[2], b.raw[3]}
	b.raw = b.raw[4:]

	return nil
}

// Size returns the complete size of the box
func (b *box) Size() uint64 {
	if b == nil {
		return 0
	}

	if b.size == 1 {
		return b.largesize
	}
	return uint64(b.size)
}

// Type returns the 4-character box type
func (b *box) Type() string {
	return b.boxtype
}

// Raw returns the raw bytes of the box
func (b *box) Raw() []byte {
	return b.raw
}

func readBoxes(buf []byte) <-chan *box {
	boxes := make(chan *box)
	r := bytes.NewReader(buf)
	go func() {
		for eof := false; !eof; {
			b, err := newBox(r)
			if err != nil {
				switch errors.Cause(err) {
				case io.EOF:
					eof = true
				}
			}

			if b != nil {
				boxes <- b
			}
		}
		close(boxes)
	}()
	return boxes
}
