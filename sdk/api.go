package sdk

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type Request struct {
	Request *http.Request
}

func newRequest(request *http.Request) *Request {
	return &Request{Request: request}
}

// InputFs assumes the request to be of type multipart/form-data and extracts the inputFs from the request.
//
// The inputFs parameter is expected to be a form file containing a zip archive.
func (r *Request) InputFs() (*InputFs, error) {
	err := r.Request.ParseMultipartForm(0)
	if err != nil {
		return nil, err
	}

	file, header, err := r.Request.FormFile("inputFs")
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(file, header.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip reader: %w", err)
	}

	return newInputFsFromZip(zipReader), nil
}

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
