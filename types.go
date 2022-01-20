package main

const (
	TypeGolangSource = "type-product-os-t-golang-source@1.1.0"
	TypeExecutable   = "type-product-os-t-executable@1.1.0"
	TypeTestRun      = "type-product-os-t-test-run@1.0.2"
)

// InputManifest describes the input to the transformer
//
// ref: https://github.com/product-os/transformer-runtime/blob/v1.5.0/lib/types/index.ts#L47
type InputManifest struct {
	Input struct {
		// ArtifactPath is the directory containing assets,
		// must be relative to the manifest
		ArtifactPath string   `json:"artifactPath"`
		Contract     Contract `json:"contract"`
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
	ArtifactPath string   `json:"artifactPath"`
	Contract     Contract `json:"contract"`
}

// Contract is a JellyFish contract
//
// ref: https://github.com/product-os/transformer-runtime/blob/v1.5.0/lib/types/index.ts#L14
// ref: https://github.com/product-os/jellyfish-types/blob/v2.0.4/lib/core/contracts/contract.ts#L11
type Contract struct {
	Type    string       `json:"type"`
	Version string       `json:"version,omitempty"`
	Name    string       `json:"name,omitempty"`
	Data    ContractData `json:"data"`
}

// ContractData is the type-specific part of Contract
//
// ref: https://github.com/product-os/transformer-runtime/blob/v1.5.0/lib/types/index.ts#L24
type ContractData struct {
	// NOTE: We compose this meta type out of the combined fields of all types
	// supported below. When adding new fields be careful about overlapping
	// definitions. All fields need to specify `omitempty`.
	//
	// TODO: a better approach might to dynamically unmarshal this at runtime
	GolangSourceData
	ExecutableData
	TestRunData
}

// GolangSourceData describes a repository, containing Go source code
//
// ref: https://github.com/product-os/t-golang-source
type GolangSourceData struct {
	// Platforms indicates which target platforms we support
	Platforms []string `json:"platforms,omitempty"`
	// Binaries lists the targets this repo can be built into,
	// they should appear under ./cmd/<binary>
	// If this is omitted we default to the repository name.
	Binaries []string `json:"binaries,omitempty"`
	// Tags lists Go build tags, which should we set for every build.
	// TODO: do we need to specify them per binary
	Tags []string `json:"tags,omitempty"`
	// DependsOn indicates system-level dependencies per distribution.
	// Currently only debian packages are supported.
	DependsOn map[string][]string `json:"dependsOn,omitempty"`
}

// ExecutableData describes a single executable output
//
// ref: https://github.com/product-os/t-executable
type ExecutableData struct {
	// Platforms indicates the target platform this executable supports
	Platform string `json:"platform,omitempty"`
	// Filename of the executable
	Filename string `json:"filename,omitempty"`
	// Version of the executable
	Version string `json:"version,omitempty"`
	// DependsOn indicates system-level dependencies per distribution
	DependsOn map[string][]string `json:"dependsOn,omitempty"`
}

// TestRunData describes the result of a test run
//
// ref: https://github.com/product-os/t-test-run
type TestRunData struct {
	// Success indicates the success status of the test run
	Success bool          `json:"success,omitempty"`
	Suites  []SuiteResult `json:"suiteResults,omitempty"`
}

type SuiteResult struct {
	Name           string   `json:"suiteName,omitempty"`
	Success        bool     `json:"suiteSuccess,omitempty"`
	UnmatchedFiles []string `json:"unmatchedFiles,omitempty"`
}
