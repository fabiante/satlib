package main

import (
	_ "embed"
	"errors"
	"log/slog"
	"os/exec"
	"path/filepath"

	"github.com/fabiante/satlib/sdk"
	"github.com/spf13/afero"
)

//go:embed openapi.yml
var openApiSpec []byte

func main() {
	sdk.Setup(
		sdk.WithOpenApiSpec(openApiSpec),
	).Run(Handle)
}

func Handle(ctx *sdk.Context) error {
	request, response := ctx.Http()

	// Setup workdir
	workdir, err := sdk.WithWorkdir(request)
	if err != nil {
		return err
	}

	inputFs, outputFs := workdir.Input, workdir.Output

	// Remove the temporary fs after the handler has finished
	defer func() {
		if err := workdir.Unlink(); err != nil {
			slog.Error("failed to unlink workdir", "err", err)
		}
	}()

	// Assumption: The input directory contains a single file called "input.pdf"
	// Let's assert that that is true
	if exists, err := afero.Exists(inputFs, "input.pdf"); err != nil {
		return err
	} else if !exists {
		return errors.New("expected input.pdf to be present")
	}

	// Run command to render PDF into images
	cmd := newCmd(sdk.FullFsPath(inputFs, "input.pdf"), sdk.FullFsPath(outputFs, ""))
	slog.Info("running cmd", "cmd", cmd.String())
	if err := cmd.Run(); err != nil {
		return err
	}

	// Zip output directory and write to response
	return response.WriteZipFs(outputFs)
}

func newCmd(inputPath string, outBasePath string) *exec.Cmd {
	program := "pdftoppm"
	args := []string{
		inputPath,
		filepath.Join(outBasePath, "page"),
		"-png",
	}
	cmd := exec.Command(program, args...)
	return cmd
}
