package sdk

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

type NestedFs = *afero.BasePathFs

func FullFsPath(fs NestedFs, name string) string {
	return afero.FullBaseFsPath(fs, name)
}

func NewTempFs() (NestedFs, func() error) {
	root := filepath.Join(os.TempDir(), uuid.New().String())
	fs := afero.NewOsFs()
	unlink := func() error {
		return fs.RemoveAll(root)
	}
	return NewNestedFs(fs, root), unlink
}

func NewNestedFs(parent afero.Fs, path string) NestedFs {
	fs := afero.NewBasePathFs(parent, path)
	if err := fs.MkdirAll("/", 0755); err != nil {
		panic(err)
	}

	return fs.(NestedFs)
}
