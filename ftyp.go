package bmff

import (
	"encoding/binary"

	"github.com/davecgh/go-spew/spew"
)

// FileType is short-formed as 'ftyp'.
// Files written to this version of this specification must contain a file‐type box.
// For compatibility with an earlier version of this specification, files may be
// conformant to this specification and not contain a file-type box.
// Files with no file‐type box should be read as if they contained an FTYP box
// with Major_brand='mp41', minor_version=0, and the single compatible brand 'mp41'.
type FileType struct {
	*box
	MajorBrand       string
	MinorVersion     int
	CompatibleBrands []string
}

func (b *FileType) parse() error {
	b.MajorBrand, b.MinorVersion = string(b.raw[0:4]), int(binary.BigEndian.Uint32(b.raw[4:8]))
	if l := len(b.raw); l > 8 {
		for i := 8; i < l; i += 4 {
			b.CompatibleBrands = append(b.CompatibleBrands, string(b.raw[i:i+4]))
		}
	}

	b.raw = nil
	spew.Dump(b)
	return nil
}
