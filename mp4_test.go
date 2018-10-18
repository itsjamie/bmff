package bmff

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func readerFromFixture(t *testing.T, path string) io.Reader {
	t.Helper()

	f, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		t.Fatalf("failed to open %s file: %v", path, err)
	}

	return f
}

func TestParse(t *testing.T) {
	type args struct {
		src io.Reader
	}

	tests := []struct {
		name    string
		args    args
		want    File
		wantErr bool
	}{
		{
			name: "Test box parsing",
			args: args{
				src: readerFromFixture(t, filepath.Join("testdata", "01_simple.mp4")),
			},
			want: File{
				Movie: &Movie{
					box: &box{
						boxtype: "ftyp",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// spew.Dump(f.Movie.Tracks[0].Media)
		})
	}
}
