package sdk

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type Response struct {
	Writer http.ResponseWriter
}

func newResponse(writer http.ResponseWriter) *Response {
	return &Response{Writer: writer}
}

func (r *Response) WriteJSON(status int, data any) error {
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(status)

	return json.NewEncoder(r.Writer).Encode(data)

}

func (r *Response) WriteZipFs(sourceFs afero.Fs) error {
	r.Writer.Header().Set("Content-Type", "application/zip")
	r.Writer.Header().Set("Content-Disposition", "attachment; filename=output.zip")

	zipWriter := zip.NewWriter(r.Writer)
	defer func() {
		err := zipWriter.Close()
		if err != nil {
			panic(err)
		}
	}()

	copyFile := func(path string, info fs.FileInfo) error {
		read, err := sourceFs.OpenFile(path, os.O_RDONLY, info.Mode())
		if err != nil {
			return err
		}
		defer func() {
			err := read.Close()
			if err != nil {
				panic(err)
			}
		}()

		write, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		if _, err := io.Copy(write, read); err != nil {
			return err
		}

		return nil
	}

	return afero.Walk(sourceFs, "/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return copyFile(path, info)
	})
}
