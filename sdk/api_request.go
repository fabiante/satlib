package sdk

import (
	"archive/zip"
	"fmt"
	"net/http"
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
