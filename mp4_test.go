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
		// {
		// 	name: "01_simple.mp4",
		// 	args: args{
		// 		src: readerFromFixture(t, filepath.Join("testdata", "01_simple.mp4")),
		// 	},
		// 	want: File{
		// 		Movie: &Movie{
		// 			box: &box{
		// 				boxtype: "ftyp",
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "02_dref_edts_img.mp4",
		// 	args: args{
		// 		src: readerFromFixture(t, filepath.Join("testdata", "02_dref_edts_img.mp4")),
		// 	},
		// 	want: File{},
		// },
		// {
		// 	name: "03_hinted.mp4",
		// 	args: args{
		// 		src: readerFromFixture(t, filepath.Join("testdata", "03_hinted.mp4")),
		// 	},
		// 	want: File{},
		// },
		// {
		// 	name: "04_bifs_video.mp4",
		// 	args: args{
		// 		src: readerFromFixture(t, filepath.Join("testdata", "04_bifs_video.mp4")),
		// 	},
		// 	want: File{},
		// },
		// {
		// 	name: "05_bifs_video_protected_v2.mp4",
		// 	args: args{
		// 		src: readerFromFixture(t, filepath.Join("testdata", "05_bifs_video_protected_v2.mp4")),
		// 	},
		// 	want: File{},
		// },
		// {
		// 	name: "06_bifs.mp4",
		// 	args: args{
		// 		src: readerFromFixture(t, filepath.Join("testdata", "06_bifs.mp4")),
		// 	},
		// 	want: File{},
		// },
		{
			name: "07_bifs_sprite.mp4",
			args: args{
				src: readerFromFixture(t, filepath.Join("testdata", "07_bifs_sprite.mp4")),
			},
			want: File{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// spew.Dump(f)
		})
	}
}
