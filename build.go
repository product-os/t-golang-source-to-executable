package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/product-os/t-golang-source-to-executable/pkg/shell"
)

type BuildOpts struct {
	Bin       string
	OutputDir string
	Version   string
	Tags      []string
	Hack      hack
}

func build(workdir string, debug bool, opts BuildOpts) error {
	var buildEnv []string

	// handle non-modules
	workdir, buildEnv, err := gopathFix(workdir, opts, buildEnv)
	if err != nil {
		return fmt.Errorf("error setting up GOPATH: %w", err)
	}

	buildArgs := []string{
		"build",
		"-o", filepath.Join(opts.OutputDir, opts.Bin),

		// TODO: is standardizing on a version pkg a good idea?
		//
		// moby uses "dockerversion.{Version,Revision}"
		// buildkit uses "version.{Version,Revision}"
		// containerd uses "version.{Version,Revision}"
		//
		// https://pkg.go.dev/runtime/debug#ReadBuildInfo
		// actually supports getting the embedded module version since go 1.12
		// we could rely on projects using that to set version... instead of embedding
		// it like this
		"-ldflags", fmt.Sprintf("-X version.Version=%s", opts.Version),

		// TODO: is the context always a git repo / can I get the git commit
		// with https://github.com/golang/go/issues/37475 (go 1.18)
		// git revision information will be automatically added by the go tool,
		// https://github.com/carlmjohnson/versioninfo is an easy way to set those
		// values without settings `-ldflags "-X ..."` from an external wrapper script
		// "-ldflags", fmt.Sprintf("-X version.Revision=%s", revision?),
	}

	if debug {
		buildArgs = append(buildArgs, "-x")
	}

	if len(opts.Tags) > 0 {
		buildArgs = append(buildArgs,
			"-tags", strings.Join(opts.Tags, ","))
	}

	// NOTE: not prepending the `./` gives just `cmd/<name>`
	// which leads go to look for a package in GOPATH :/
	buildArgs = append(buildArgs, "./"+filepath.Join(".", "cmd", opts.Bin))

	_, err = shell.Run("go", buildArgs, nil, os.Stdout, os.Stderr,
		shell.WithDir(workdir),
		shell.WithEnv(buildEnv))

	return err
}

func gopathFix(workdir string, opts BuildOpts, env []string) (string, []string, error) {
	_, err := os.Stat(filepath.Join(workdir, "go.mod"))
	// is go module
	if err == nil {
		log.Println("build mode = module")
		return workdir, env, nil
	}
	// handle stat error
	if !errors.Is(err, os.ErrNotExist) {
		return workdir, env, err
	}

	// needs gopath fix
	log.Println("build mode = gopath")
	// disable go modules
	env = append(env, "GO111MODULE=off")

	// construct gopath location
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		return "", nil, errors.New("GOPATH undefined")
	}
	// set $GOPATH/src/<bin>
	module := opts.Bin
	// HACK: enable setting $GOPATH/src/<hack.module>
	if opts.Hack.Module != "" {
		module = opts.Hack.Module
	}
	gopath = filepath.Join(gopath, "src", module)
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
