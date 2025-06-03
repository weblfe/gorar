package detect_test

import (
	"context"
	"github.com/mholt/archives"
	"github.com/weblfe/gorar/pkg/detect"
	"io/fs"
	"testing"
)

type Case struct {
	filePath string
	want     string
	fss      fs.FS
	create   func(fs fs.FS, filePath string)
}

func newFs() fs.FS {
	ctx := context.Background()
	fss, _ := archives.FileSystem(ctx, "./", nil)
	return fss
}

func autoCreateArchives(items []*Case) {
	var fss = newFs()
	for _, it := range items {
		it.fss = fss
		if it.create == nil {
			it.create = createArchive(it.want)
		}
		it.create(fss, it.filePath)
	}
}

func createArchive(format string) func(fs fs.FS, filePath string) {
	switch format {
	case `zip`:

	case `xz`:
	case `gz`:

	}
	return func(fs fs.FS, filePath string) {

	}
}

func Test_DetectType(t *testing.T) {
	var (
		cases = []*Case{
			{
				filePath: "testdata/test.zip",
				want:     "zip",
			},
			{
				filePath: "testdata/test.tar.gz",
				want:     "gz",
			}, {
				filePath: "testdata/test.tar.xz",
				want:     "xz",
			},
			{
				filePath: "testdata/test.tar",
				want:     "tar",
			},
			{
				filePath: "testdata/test.gz",
				want:     "gz",
			},
			{
				filePath: "testdata/test.xz",
				want:     "xz",
			},
			{
				filePath: "testdata/test.rar",
				want:     "rar",
			},
			{
				filePath: "testdata/test.7z",
				want:     "7z",
			},
		}
	)
	autoCreateArchives(cases)
	for _, c := range cases {
		if got, _ := detect.Detect(c.filePath, c.fss); got != c.want {
			t.Errorf("Detect(%s) = %s, want %s", c.filePath, got, c.want)
		}
	}
}
