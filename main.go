package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/product-os/t-golang-source-to-executable/pkg/packages"
	"github.com/product-os/t-golang-source-to-executable/pkg/shell"
)

const (
	OUTPUT_ARTIFACT_DIRNAME = "artifacts"
)

var (
	inputPath    string
	outputPath   string
	mode         string
	artifactPath string
	debug        bool
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Flags() | log.Lshortfile)

	flag.StringVar(&inputPath, "input", "", "input contract path")
	flag.StringVar(&outputPath, "output", "", "output contract path")
	flag.StringVar(&mode, "mode", "build", "what do? [build or test]")
	flag.StringVar(&artifactPath, "artifactPath", OUTPUT_ARTIFACT_DIRNAME, "artifact path for output assets")
	flag.BoolVar(&debug, "debug", false, "be verbose")
	flag.Parse()

	flag.CommandLine.VisitAll(func(fl *flag.Flag) {
		// override with envvar if set
		if val, ok := os.LookupEnv(strings.ToUpper(fl.Name)); ok {
			if err := fl.Value.Set(val); err != nil {
				panic(err)
			}
		}
		// print values
		log.Printf("%s = %q", fl.Name, fl.Value.String())
	})

	if debug {
		var err error
		_, err = shell.Run("go", []string{"env"}, nil, os.Stdout, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	var (
		input  InputManifest
		output OutputManifest
	)

	// load input contract
	inputFd, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(inputFd).Decode(&input); err != nil {
		log.Fatal(err)
	}

	if debug {
		log.Printf("INPUT:\n%#+v", input)
	}

	inputArtifactPath := filepath.Join(filepath.Dir(inputPath), input.Input.ArtifactPath)
	outputArtifactPath := filepath.Join(filepath.Dir(outputPath), artifactPath)
	if err := os.MkdirAll(outputArtifactPath, os.ModeDir|os.ModePerm); err != nil {
		log.Fatalf("creating output artifact path %s: %v", outputArtifactPath, err)
	}

	if err := os.Chdir(inputArtifactPath); err != nil {
		log.Fatal(err)
	}

	if err := setupTaskEnvironment(input.Input.Contract.Data.GolangSourceData); err != nil {
		log.Fatal(err)
	}

	switch mode {
	case "build":
		var binaries []string
		if input.Input.Contract.Data.Binaries != nil {
			binaries = input.Input.Contract.Data.Binaries
		} else {
			// fallback to a binary with the same name as the repo
			binaries = []string{input.Input.Contract.Name}
		}
		for _, bin := range binaries {
			if err := build(inputArtifactPath, debug, BuildOpts{
				Name:      bin,
				Version:   input.Input.Contract.Version,
				Tags:      input.Input.Contract.Data.GolangSourceData.Tags,
				OutputDir: outputArtifactPath,
			}); err != nil {
				log.Fatalf("build failed: %v", err)
			}
			output.Results = append(output.Results, Result{
				ArtifactPath: artifactPath,
				Contract: Contract{
					Type: TypeExecutable,
					Data: ContractData{ExecutableData: ExecutableData{
						// NOTE: we set platform to the native platform of the go runtime here
						// this means we completely disregard the fact that we could actually
						// do cross-compilation.
						// We have the target platform in `input.Input.Contract.Data.GolangSourceData.Platforms`
						Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
						Filename:  bin,
						Version:   input.Input.Contract.Version,
						DependsOn: input.Input.Contract.Data.GolangSourceData.DependsOn,
					}},
				},
			})
		}

	case "test":
		unitSuites, unitErr := test(inputArtifactPath, debug, TestOpts{
			Name:        input.Input.Contract.Name,
			Tags:        input.Input.Contract.Data.GolangSourceData.Tags,
			Integration: false,
		})
		if unitSuites == nil && unitErr != nil {
			log.Fatalf("unit tests failed: %v", unitErr)
		}
		integrationSuites, integrationErr := test(inputArtifactPath, debug, TestOpts{
			Name:        input.Input.Contract.Name,
			Tags:        input.Input.Contract.Data.GolangSourceData.Tags,
			Integration: true,
		})
		if integrationSuites == nil && integrationErr != nil {
			log.Fatalf("integration tests failed: %v", integrationErr)
		}
		output.Results = append(output.Results, Result{
			ArtifactPath: artifactPath,
			Contract: Contract{
				Type: TypeTestRun,
				Data: ContractData{TestRunData: TestRunData{
					Success: unitErr == nil && integrationErr == nil,
					Suites:  append(unitSuites, integrationSuites...),
				}},
			},
		})

	default:
		log.Fatalf("unknown mode %q", mode)

	}

	// write output contract
	outputFd, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("creating output manifest %s: %v", outputPath, err)
	}
	if err := json.NewEncoder(outputFd).Encode(output); err != nil {
		log.Fatalf("writing output manifest %s: %v", outputPath, err)
	}
}

func setupTaskEnvironment(data GolangSourceData) error {
	if data.DependsOn != nil {
		log.Println("fetching dependencies")
		for distro, pkgs := range data.DependsOn {
			if err := packages.Install(distro, pkgs...); err != nil {
				return fmt.Errorf("error installing packages for %s: %w", distro, err)
			}
		}
	}
	return nil
}
