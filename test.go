package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/product-os/t-golang-source-to-executable/pkg/shell"
)

type TestOpts struct {
	Name     string
	Contract *GolangSourceData
}

func test(workdir string, debug bool, opts TestOpts) (suites []SuiteResult, err error) {
	var testEnv []string

	// handle non-modules
	workdir, testEnv, err = gopathFix(workdir, opts.Name, opts.Contract, testEnv)
	if err != nil {
		return nil, fmt.Errorf("error setting up GOPATH: %w", err)
	}

	testReportFd, err := ioutil.TempFile(os.TempDir(), "tf-golang-test-report-XXXX.json")
	if err != nil {
		return nil, err
	}

	// gotestsum args
	testArgs := []string{
		"--format", "standard-verbose",
		"--jsonfile", testReportFd.Name(),
		"--",
	}

	// go test args
	testArgs = append(testArgs, "-timeout=5m", "-cover", "-covermode=atomic")

	// go build args
	var buildArgs []string
	if len(opts.Contract.Tags) > 0 {
		buildArgs = append(buildArgs,
			"-tags", strings.Join(opts.Contract.Tags, ","))
	}
	testArgs = append(testArgs, buildArgs...)

	pkgList, err := listPackages(workdir, buildArgs, testEnv)
	if err != nil {
		return nil, err
	}

	// skip if there's nothing to test
	if len(pkgList) < 1 {
		log.Printf("no tests to run")
		return nil, err
	}

	testArgs = append(testArgs, pkgList...)

	_, err = shell.Run("gotestsum", testArgs, nil, os.Stdout, os.Stderr,
		shell.WithDir(workdir),
		shell.WithEnv(testEnv))

	scanner := bufio.NewScanner(testReportFd)
	for scanner.Scan() {
		var report map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &report); err != nil {
			return nil, fmt.Errorf("error parsing gotestsum report %s: %w", testReportFd.Name(), err)
		}
		if report["Test"] != nil ||
			(report["Action"] != "pass" &&
				report["Action"] != "fail") {
			continue
		}
		suites = append(suites, SuiteResult{
			Success: (report["Action"] == "pass"),
			Name:    fmt.Sprintf("%v", report["Package"]),
		})
	}

	// TODO: should we copy the test-report json to artifacts?
	return suites, err
}

func listPackages(workdir string, buildArgs, env []string) (pkgList []string, err error) {
	var stdout bytes.Buffer
	listArgs := []string{"list"}
	listArgs = append(listArgs, buildArgs...)
	listArgs = append(listArgs, workdir+"/...")
	_, err = shell.Run("go", listArgs, nil, &stdout, nil,
		shell.WithDir(workdir),
		shell.WithEnv(env))
	if err != nil {
		return nil, err
	}
	for scanner := bufio.NewScanner(&stdout); scanner.Scan(); {
		pkgList = append(pkgList, scanner.Text())
	}
	return pkgList, nil
}
