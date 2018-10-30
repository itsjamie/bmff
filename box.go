package bmff

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

const boxHeaderSize = 8

func newBox(src io.Reader) (*box, error) {
	b := &box{}
	if err := b.decode(src); err != nil {
		return nil, err
	}
	return b, nil
}

type box struct {
	boxtype   string
	size      uint32
	largesize uint64
	raw       []byte
}

func (b *box) decode(src io.Reader) error {
	header := make([]byte, 8)
	if _, err := io.ReadFull(src, header); err != nil {
		return errors.Wrapf(err, "error reading box header: %x", header)
	}

	b.size = binary.BigEndian.Uint32(header[0:4])
	b.boxtype = string(header[4:8])

	if b.size == 1 {
		if err := binary.Read(src, binary.BigEndian, &b.largesize); err != nil {
			return errors.Wrap(err, "failure to read largesize field")
		}
	}

	if b.size == 0 {
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

	return nil
}

func (b *box) Size() uint64 {
	if b == nil {
		return 0
	}

	if b.size == 1 {
		return b.largesize
	}
	return uint64(b.size)
}

func (b *box) Type() string {
	return b.boxtype
}

func (b *box) Raw() []byte {
	return b.raw
}

type fullbox struct {
	*box
	version uint8
	flags   uint32
}

func (b *fullbox) decode() error {
	if len(b.raw) < 4 {
		return errors.Errorf("failed decode of fullbox type %s, missing data", b.box.boxtype)
	}

	b.version = uint8(b.raw[0])
	b.flags = binary.BigEndian.Uint32([]byte{0x00, b.raw[1], b.raw[2], b.raw[3]})
	b.raw = b.raw[4:]

	return nil
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
				default:
					log.Fatal(err)
				}
			}

			if b != nil {
				log.Println("newBox - " + b.boxtype)
				boxes <- b
			}
		}
		close(boxes)
	}()
	return boxes
}
