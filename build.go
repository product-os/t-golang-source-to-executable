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
	Name      string
	OutputDir string
	Version   string
	Tags      []string
}

func build(workdir string, debug bool, opts BuildOpts) error {
	var buildEnv []string

	// handle non-modules
	workdir, buildEnv, err := gopathFix(workdir, opts.Name, buildEnv)
	if err != nil {
		return fmt.Errorf("error setting up GOPATH: %w", err)
	}

	buildArgs := []string{
		"build",
		"-o", filepath.Join(opts.OutputDir, opts.Name),
		// TODO: investigate https://github.com/golang/go/issues/37475
		//                   https:// github.com/golang/go/commit/76cef81bcff371c88d277f17c712ecf22b8c83e7
		//                   (needs go1.18)
		"-ldflags", fmt.Sprintf("-X version.Version=%s", opts.Version),
	}

	if debug {
		buildArgs = append(buildArgs, "-x")
	}

	if len(opts.Tags) > 0 {
		buildArgs = append(buildArgs,
			"-tags", strings.Join(opts.Tags, " "))
	}

	// NOTE: not prepending the `./` gives just `cmd/<name>`
	// which leads go to look for a package in GOPATH :/
	buildArgs = append(buildArgs, "./"+filepath.Join(".", "cmd", opts.Name))

	_, err = shell.Run("go", buildArgs, nil, os.Stdout, os.Stderr,
		shell.WithDir(workdir),
		shell.WithEnv(buildEnv))

	return err
}

func gopathFix(workdir string, name string, env []string) (string, []string, error) {
	_, err := os.Stat(filepath.Join(workdir, "go.mod"))
	if err == nil {
		return workdir, env, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return workdir, env, err
	}
	env = append(env, "GO111MODULE=off")
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		return "", nil, errors.New("GOPATH undefined")
	}
	gopath = filepath.Join(gopath, "src", name)
	// check if it exists and exit early
	if _, err = os.Stat(gopath); err == nil {
		return gopath, env, nil
	}
	// create "fake" gopath entry
	if err := os.MkdirAll(filepath.Dir(gopath), os.ModeDir|os.ModePerm); err != nil {
		return "", nil, err
	}
	if err := os.Symlink(workdir, gopath); err != nil {
		return "", nil, err
	}
	return gopath, env, nil
}
