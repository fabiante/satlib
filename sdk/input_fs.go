package sdk

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type InputFs struct {
	Zip *zip.Reader
}

func newInputFsFromZip(zip *zip.Reader) *InputFs {
	return &InputFs{Zip: zip}
}

func (fs *InputFs) CopyTo(target afero.Fs) error {
	copyFile := func(f *zip.File) error {
		fr, err := f.Open()
		if err != nil {
			return err
		}
		defer func(fr io.ReadCloser) {
			err := fr.Close()
			if err != nil {
				panic(err)
			}
		}(fr)

		if f.Mode().IsDir() {
			if err := target.MkdirAll(f.Name, f.Mode()); err != nil {
				return err
			}
		} else {
			if err := target.MkdirAll(filepath.Dir(f.Name), f.Mode()); err != nil {
				return err
			}

			tw, err := target.OpenFile(f.Name, os.O_CREATE|os.O_WRONLY, f.Mode())
			if err != nil {
				return err
			}
			defer func(tw afero.File) {
				err := tw.Close()
				if err != nil {
					panic(err)
				}
			}(tw)

			_, err = io.Copy(tw, fr)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// Iterate over all ZIP files and extract them to the given directory.
	for _, f := range fs.Zip.File {
		if err := copyFile(f); err != nil {
			return err
		}
	}

	return nil
}
