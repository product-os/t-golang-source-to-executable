package main

import (
	"encoding/json"
	"fmt"
)

const (
	TypeGolangSource = "type-product-os-t-golang-source@1.1.3"
	TypeExecutable   = "type-product-os-t-executable@1.1.0"
	TypeTestRun      = "type-product-os-t-test-run@1.1.0"
)

// InputManifest describes the input to the transformer
//
// ref: https://github.com/product-os/transformer-runtime/blob/v1.5.0/lib/types/index.ts#L47
type InputManifest struct {
	Input struct {
		// ArtifactPath is the directory containing assets,
		// must be relative to the manifest
		ArtifactPath string `json:"artifactPath"`
		// Contract is the contract describing the input
		Contract Contract `json:"contract"`
	} `json:"input"`
}

// OutputManifest describes the output of the transformer
//
// ref: https://github.com/product-os/transformer-runtime/blob/v1.5.0/lib/types/index.ts#L61
type OutputManifest struct {
	Results []Result `json:"results"`
}

type Result struct {
	// ArtifactPath is the directory containing assets,
	// must be relative to the manifest
	ArtifactPath string `json:"artifactPath,omitempty"`
	// Contract is the contract describing the result,
	Contract Contract `json:"contract"`
}

// Contract is a JellyFish contract
//
// ref: https://github.com/product-os/transformer-runtime/blob/v1.5.0/lib/types/index.ts#L14
// ref: https://github.com/product-os/jellyfish-types/blob/v2.0.4/lib/core/contracts/contract.ts#L11
type contractJSON struct {
	Type    string          `json:"type"`
	Version string          `json:"version,omitempty"`
	Name    string          `json:"name,omitempty"`
	Data    json.RawMessage `json:"data"`
}
type Contract struct {
	Type    string
	Version string
	Name    string
	Data    ContractData
}

func (c *Contract) UnmarshalJSON(data []byte) (err error) {
	var intermediate contractJSON
	err = json.Unmarshal(data, &intermediate)
	if err != nil {
		return err
	}
	switch intermediate.Type {
	case TypeGolangSource:
		c.Data.GolangSourceData = new(GolangSourceData)
		err = json.Unmarshal(intermediate.Data, c.Data.GolangSourceData)
	case TypeExecutable:
		c.Data.ExecutableData = new(ExecutableData)
		err = json.Unmarshal(intermediate.Data, c.Data.ExecutableData)
	case TypeTestRun:
		c.Data.TestRunData = new(TestRunData)
		err = json.Unmarshal(intermediate.Data, c.Data.TestRunData)
	default:
		return fmt.Errorf("unknown contract type %q", intermediate.Type)
	}
	if err != nil {
		return err
	}
	c.Type = intermediate.Type
	c.Version = intermediate.Version
	c.Name = intermediate.Name
	return err
}

func (c *Contract) MarshalJSON() ([]byte, error) {
	var (
		intermediate contractJSON
		err          error
	)
	switch c.Type {
	case TypeGolangSource:
		intermediate.Data, err = json.Marshal(c.Data.GolangSourceData)
	case TypeExecutable:
		intermediate.Data, err = json.Marshal(c.Data.ExecutableData)
	case TypeTestRun:
		intermediate.Data, err = json.Marshal(c.Data.TestRunData)
	default:
		return nil, fmt.Errorf("unknown contract type %q", intermediate.Type)
	}
	if err != nil {
		return nil, err
	}
	intermediate.Type = c.Type
	intermediate.Version = c.Version
	intermediate.Name = c.Name
	return json.Marshal(intermediate)
}

// ContractData is the type-specific part of Contract
//
// ref: https://github.com/product-os/transformer-runtime/blob/v1.5.0/lib/types/index.ts#L24
type ContractData struct {
	GolangSourceData *GolangSourceData
	ExecutableData   *ExecutableData
	TestRunData      *TestRunData
}

// GolangSourceData describes a repository, containing Go source code
//
// ref: https://github.com/product-os/t-golang-source
type GolangSourceData struct {
	// Platforms indicates which target platforms we support
	Platforms []string `json:"platforms"`
	// Binaries lists the targets this repo can be built into,
	// they should appear under ./cmd/<binary>
	// If this is omitted we default to the repository name.
	Binaries []string `json:"binaries,omitempty"`
	// Tags lists Go build tags, which should we set for every build.
	// TODO: do we need to specify them per binary
	Tags []string `json:"tags,omitempty"`
	// DependsOn indicates system-level dependencies per distribution.
	// TODO: currently only debian packages are supported.
	DependsOn map[string][]string `json:"dependsOn,omitempty"`
}

// ExecutableData describes a single executable output
//
// ref: https://github.com/product-os/t-executable
type ExecutableData struct {
	// Platforms indicates the target platform this executable supports
	Platform string `json:"platform"`
	// Filename of the executable
	Filename string `json:"filename"`
	// Version of the executable
	Version string `json:"version"`
	// DependsOn indicates system-level dependencies per distribution
	DependsOn map[string][]string `json:"dependsOn,omitempty"`
}

// TestRunData describes the result of a test run
//
// ref: https://github.com/product-os/t-test-run
type TestRunData struct {
	// Success indicates the success status of the test run
	Success bool          `json:"success"`
	Suites  []SuiteResult `json:"suiteResults,omitempty"`
}

type SuiteResult struct {
	Name           string   `json:"suiteName"`
	Success        bool     `json:"suiteSuccess"`
	UnmatchedFiles []string `json:"unmatchedFiles,omitempty"`
}
