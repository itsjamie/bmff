package bmff

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

const boxHeaderSize = 8

type Box interface {
	Type() string
	Size() uint64
	Raw() []byte
}

type File struct {
	Ftyp    *FtypBox
	Moov    *MoovBox
	Unknown []Box
}

// Box is defined as "object‐oriented building block defined by a unique type identifier and length".
// Used in MP4 Containers and referenced as an "atom" in some specifications including the first definition of MP4.
//
// Boxes with an unrecognized type should be ignored and skipped.
type box struct {
	boxtype   string
	size      uint32
	largesize uint64
	raw       []byte
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

func NewBox(src io.Reader) (*box, error) {
	buf := make([]byte, boxHeaderSize)
	_, err := io.ReadAtLeast(src, buf, 8)
	if err != nil {
		return nil, errors.Wrap(err, "error reading buffer header")
	}
	s := binary.BigEndian.Uint32(buf[0:4])
	b := &box{
		boxtype: string(buf[4:8]),
		size:    s,
		raw:     make([]byte, s-boxHeaderSize),
	}

	_, err = io.ReadFull(src, b.raw)
	if err != nil {
		return b, errors.Wrap(err, "error reading box data")
	}

	return b, nil
}

type FtypBox struct {
	*box
	MajorBrand       string
	MinorVersion     int
	CompatibleBrands []string
}

func (b *FtypBox) parse() error {
	b.MajorBrand, b.MinorVersion = string(b.raw[0:4]), int(binary.BigEndian.Uint32(b.raw[4:8]))
	if l := len(b.raw); l > 8 {
		for i := 8; i < l; i += 4 {
			b.CompatibleBrands = append(b.CompatibleBrands, string(b.raw[i:i+4]))
		}
	}
	return nil
}

type MoovBox struct {
	*box
	MovieHeader *MvhdBox
	TrackBoxes  []*TrakBox
	Iods        *IodsBox
	// Meta *MetaBox

}

func (b *MoovBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		if subBox == nil {
			return nil
		}

		switch subBox.boxtype {
		case "mvhd":
			b.MovieHeader = &MvhdBox{box: subBox}
			if err := b.MovieHeader.parse(); err != nil {
				return err
			}
		case "iods":
			b.Iods = &IodsBox{box: subBox}
			if err := b.Iods.parse(); err != nil {
				return err
			}
		case "trak":
			trak := &TrakBox{box: subBox}
			if err := trak.parse(); err != nil {
				return err
			}

			b.TrackBoxes = append(b.TrackBoxes, trak)
		default:
			fmt.Printf("Unknown Moov SubType: %s\n", subBox.Type())
		}
	}

	return nil
}

type TrakBox struct {
	*box
	Header           *TkhdBox
	Reference        *TrefBox
	MediaInformation *MdiaBox
}

func (b *TrakBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		if subBox == nil {
			break
		}

		switch subBox.boxtype {
		case "tkhd":
			header := &TkhdBox{box: subBox}
			if err := header.parse(); err != nil {
				return err
			}
			b.Header = header
		case "mdia":
			mdia := &MdiaBox{box: subBox}
			if err := mdia.parse(); err != nil {
				return err
			}
			b.MediaInformation = mdia
		case "tref":
			tref := &TrefBox{box: subBox}
			if err := tref.parse(); err != nil {
				return err
			}
			b.Reference = tref
		default:
			fmt.Printf("Unknown Trak SubType: %s\n", subBox.Type())
		}
	}

	return nil
}

type MdiaBox struct {
	*box
	Header *MdhdBox
}

func (b *MdiaBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		if subBox == nil {
			break
		}

		switch subBox.boxtype {
		case "mdhd":
			mdhd := MdhdBox{box: subBox}
			if err := mdhd.parse(); err != nil {
				return err
			}
			b.Header = &mdhd
		case "hdlr":
			fmt.Printf("Unhandled Mdia SubType: %s\n", subBox.Type())
		case "minf":
			fmt.Printf("Unhandled Mdia SubType: %s\n", subBox.Type())
		default:
			fmt.Printf("Unknown Mdia SubType: %s\n", subBox.Type())
		}
	}

	return nil
}

type MdhdBox struct {
	*box
	version          uint8
	flags            [3]byte
	CreationTime     uint64
	ModificationTime uint64
	TimeScale        uint32
	Duration         uint64
}

func (b *MdhdBox) parse() error {
	// spew.Dump(b.raw)

	var offset int
	if b.version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(b.raw[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(b.raw[8:12]))
		b.TimeScale = binary.BigEndian.Uint32(b.raw[12:16])
		b.Duration = uint64(binary.BigEndian.Uint32(b.raw[16:20]))
		offset = 20
	} else if b.version == 1 {
		b.CreationTime = binary.BigEndian.Uint64(b.raw[4:12])
		b.ModificationTime = binary.BigEndian.Uint64(b.raw[12:20])
		b.TimeScale = binary.BigEndian.Uint32(b.raw[20:24])
		b.Duration = binary.BigEndian.Uint64(b.raw[24:32])
		offset = 32
	}

	lang := binary.BigEndian.Uint16(b.raw[offset : offset+2]) // first bit is padding
	// spew.Dump(fmt.Sprintf("%08b", b.raw[offset:offset+2]))
	// spew.Dump(fmt.Sprintf("%08b", lang<<3))
	langcode := []byte{
		uint8(lang>>10&0x1F) + 0x60,
		uint8(lang>>15&0x1F) + 0x60, // wrong
	}
	// spew.Dump(langcode)

	// ltr1 := uint8((lang & 0xF800) >> 3)
	// ltr1 := (lang[0] & 0x7C) >> 2
	// ltr3 := (lang[3] & )
	// spew.Dump(string((uint8(ltr1) + 0x60)))
	// spew.Dump(fmt.Sprintf("%08b %08b", ltr1, 0x7C))
	// ltr2 := (lang[0] <<  )

	return nil
}

type MvhdBox struct {
	*box
	version          uint8
	flags            [3]byte
	CreationTime     uint64
	ModificationTime uint64
	TimeScale        uint32
	Duration         uint64
	NextTrackID      uint32
	Rate             Uint16_16
	Volume           Uint8_8
	Reserved         []byte
	Matrix           [9]int32
	Predefined       []byte
}

func (b *MvhdBox) parse() error {
	b.version = b.raw[0]
	b.flags = [3]byte{b.raw[1], b.raw[2], b.raw[3]}
	var offset int
	if b.version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(b.raw[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(b.raw[8:12]))
		b.TimeScale = binary.BigEndian.Uint32(b.raw[12:16])
		b.Duration = uint64(binary.BigEndian.Uint32(b.raw[16:20]))
		offset = 20
	} else if b.version == 1 {
		b.CreationTime = binary.BigEndian.Uint64(b.raw[4:12])
		b.ModificationTime = binary.BigEndian.Uint64(b.raw[12:20])
		b.TimeScale = binary.BigEndian.Uint32(b.raw[20:24])
		b.Duration = binary.BigEndian.Uint64(b.raw[24:32])
		offset = 32
	}

	if err := b.Rate.UnmarshalBinary(b.raw[offset : offset+4]); err != nil {
		return errors.Wrap(err, "failed to get rate")
	}

	if err := b.Volume.UnmarshalBinary(b.raw[offset+4 : offset+6]); err != nil {
		return errors.Wrap(err, "failed to get volume")
	}

	b.Reserved = b.raw[offset+6 : offset+16]
	offset += 16
	for i := 0; i < 9; i++ {
		b.Matrix[i] = int32(binary.BigEndian.Uint32(b.raw[offset+i : offset+i+4]))
	}
	offset += 36

	b.Predefined = b.raw[offset : offset+24]
	b.NextTrackID = binary.BigEndian.Uint32(b.raw[offset+24 : offset+28])
	return nil
}

type IodsBox struct {
	*box
}

func (b *IodsBox) parse() error {
	return nil
}

type TkhdBox struct {
	*box
	version          uint8
	flags            [3]byte
	CreationTime     uint64
	ModificationTime uint64
	TrackID          uint32
	Duration         uint64
	Layer            int16
	AlternateGroup   int16
	Volume           int16
	Matrix           [9]int32
	Width            Uint16_16
	Height           Uint16_16
}

func (b *TkhdBox) parse() error {
	b.version = b.raw[0]
	b.flags = [3]byte{b.raw[1], b.raw[2], b.raw[3]}
	var offset int
	if b.version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(b.raw[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(b.raw[8:12]))
		b.TrackID = binary.BigEndian.Uint32(b.raw[12:16])
		// 16:20 reserved
		b.Duration = uint64(binary.BigEndian.Uint32(b.raw[20:24]))
		offset = 24
	} else if b.version == 1 {
		b.CreationTime = uint64(binary.BigEndian.Uint64(b.raw[4:12]))
		b.ModificationTime = uint64(binary.BigEndian.Uint64(b.raw[12:20]))
		b.TrackID = binary.BigEndian.Uint32(b.raw[20:24])
		// 24:28 reserved
		b.Duration = uint64(binary.BigEndian.Uint64(b.raw[28:36]))
		offset = 36
	}
	offset += 8 // reserved bytes
	b.Layer = int16(binary.BigEndian.Uint16(b.raw[offset : offset+2]))
	b.AlternateGroup = int16(binary.BigEndian.Uint16(b.raw[offset+2 : offset+4]))
	b.Volume = int16(binary.BigEndian.Uint16(b.raw[offset+4 : offset+6]))
	offset += 8 // previous bytes + 2 reserved

	for i := 0; i < 9; i++ {
		b.Matrix[i] = int32(binary.BigEndian.Uint32(b.raw[offset+i : offset+i+4]))
	}
	offset += 36
	b.Width = Uint16_16(binary.BigEndian.Uint32(b.raw[offset : offset+4]))
	b.Height = Uint16_16(binary.BigEndian.Uint32(b.raw[offset+4 : offset+8]))
	return nil
}

type TrefBox struct {
	*box
	TypeBoxes []*TrefTypeBox
}

func (b *TrefBox) parse() error {
	for subBox := range readBoxes(b.raw) {
		if subBox == nil {
			break
		}
		t := TrefTypeBox{box: subBox}
		for i := 0; i < len(t.raw); i += 4 {
			t.TrackIDs = append(t.TrackIDs, binary.BigEndian.Uint32(t.raw[i:i+4]))

		}
		b.TypeBoxes = append(b.TypeBoxes, &t)
	}

	return nil
}

type TrefTypeBox struct {
	*box
	TrackIDs []uint32
}

func Parse(src io.Reader) (*File, error) {
	f := &File{}
	r := bufio.NewReader(src)

readloop:
	for {
		b, err := NewBox(r)
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

		switch b.boxtype {
		case "ftyp":
			fb := &FtypBox{box: b}
			if err := fb.parse(); err != nil {
				return nil, err
			}
			f.Ftyp = fb
		case "moov":
			mb := &MoovBox{box: b}
			if err := mb.parse(); err != nil {
				return nil, err
			}
			f.Moov = mb
		default:
			fmt.Printf("Unknown Type: %s\n", b.Type())
		}

	}

	return f, nil
}

func readBoxes(buf []byte) <-chan *box {
	boxes := make(chan *box)
	r := bytes.NewReader(buf)
	go func() {
		for eof := false; !eof; {
			b, err := NewBox(r)
			if err != nil {
				switch errors.Cause(err) {
				case io.EOF:
					eof = true
				}
			}

			boxes <- b
		}
		close(boxes)
	}()
	return boxes
}
