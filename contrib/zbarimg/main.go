package main

import (
	"bytes"
	"errors"
	"log/slog"
	"os/exec"

	"github.com/fabiante/satlib/sdk"
	"github.com/spf13/afero"
)

func main() {
	sdk.Run(Handle)
}

func Handle(ctx *sdk.Context) error {
	request, response := ctx.Http()

	// Setup workdir
	workdir, err := sdk.WithWorkdir(request)
	if err != nil {
		return err
	}

	inputFs, _ := workdir.Input, workdir.Output

	// Remove the temporary fs after the handler has finished
	defer func() {
		if err := workdir.Unlink(); err != nil {
			slog.Error("failed to unlink workdir", "err", err)
		}
	}()

	// Get all files in input dir
	var fileNames []string
	if files, err := afero.ReadDir(inputFs, "."); err != nil {
		return err
	} else {
		for _, file := range files {
			fileNames = append(fileNames, sdk.FullFsPath(inputFs, file.Name()))
		}
	}

	if len(fileNames) == 0 {
		return errors.New("no input files present")
	}

	// Run command to extract barcodes from all files
	cmd := newCmd(fileNames)
	outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	slog.Info("running cmd", "cmd", cmd.String())
	if err := cmd.Run(); err != nil {
		switch cmd.ProcessState.ExitCode() {
		case 1, 4:
			break
		default:
			return err
		}
	}

	return response.WriteJSON(200, map[string]any{
		"stdout": outBuf.String(),
		"stderr": errBuf.String(),
		"code":   cmd.ProcessState.ExitCode(),
	})
}

func newCmd(inputFilePaths []string) *exec.Cmd {
	program := "zbarimg"
	args := []string{
		"-q",
		"--xml",
	}
	args = append(args, inputFilePaths...)

	cmd := exec.Command(program, args...)
	return cmd
}
