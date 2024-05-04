package sdk

type Workdir struct {
	Work   NestedFs
	Input  NestedFs
	Output NestedFs

	unlink func() error
}

func (w *Workdir) Unlink() error {
	return w.unlink()
}

// WithWorkdir creates a work directory in the OS temp directory. Then it copies the
// request's input files to the work directory into the "input" directory.
func WithWorkdir(request *Request) (*Workdir, error) {
	// Get a temporary fs to work in.
	workFs, unlink := NewTempFs()

	// Create a nested fs to copy the request's input files to
	inputFs := NewNestedFs(workFs, "input")

	// Copy request's input files
	if ifs, err := request.InputFs(); err != nil {
		return nil, err
	} else {
		_ = ifs.CopyTo(inputFs)
	}

	// Create nested fs for output
	outputFs := NewNestedFs(workFs, "output")

	return &Workdir{
		Work:   workFs,
		Input:  inputFs,
		Output: outputFs,
		unlink: unlink,
	}, nil
}
