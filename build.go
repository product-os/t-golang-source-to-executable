package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/product-os/t-golang-source-to-executable/pkg/shell"
)

type BuildOpts struct {
	BinaryName string
	OutputDir  string
	Version    string
	Tags       []string
}

func build(workdir string, debug bool, opts BuildOpts) error {
	buildArgs := []string{
		"build",
		"-o", filepath.Join(opts.OutputDir, opts.BinaryName),
		"-ldflags", fmt.Sprintf("-X version.Version=%s", opts.Version),
	}

	if debug {
		buildArgs = append(buildArgs, "-x")
	}

	var buildEnv []string

	// handle non-modules
	_, err := os.Stat(filepath.Join(workdir, "go.mod"))
	if errors.Is(err, os.ErrNotExist) {
		gopath, ok := os.LookupEnv("GOPATH")
		if !ok {
			return errors.New("GOPATH undefined")
		}
		gopath = filepath.Join(gopath, "src", opts.BinaryName)
		if err := os.MkdirAll(filepath.Dir(gopath), os.ModeDir|os.ModePerm); err != nil {
			return err
		}
		if err := os.Symlink(workdir, gopath); err != nil {
			return err
		}
		workdir = gopath
		buildEnv = append(buildEnv, "GO111MODULE=off")
	}

	if len(opts.Tags) > 0 {
		buildArgs = append(buildArgs,
			"-tags", strings.Join(opts.Tags, " "))
	}

	// NOTE: not prepending the `./` gives just `cmd/<name>`
	// which leads go to look for a package in GOPATH :/
	buildArgs = append(buildArgs, "./"+filepath.Join(".", "cmd", opts.BinaryName))

	_, err = shell.Run("go", buildArgs, nil, os.Stdout, os.Stderr,
		shell.WithDir(workdir),
		shell.WithEnv(buildEnv))

	return err
}
