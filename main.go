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
		input  NInput
		output NOutput
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

	if err := os.Chdir(inputArtifactPath); err != nil {
		log.Fatal(err)
	}

	if err := setupTaskEnvironment(input.Input.Contract.Data.GolangSourceData); err != nil {
		log.Fatal(err)
	}

	switch mode {
	case "build":
		result := TransformerAsset{
			ArtifactPath: artifactPath,
			Contract: Contract{
				Type: TypeExecutable,
				Data: ContractData{ExecutableData: ExecutableData{
					// NOTE: we set platform to the native platform of the go runtime here
					// this means we completely disregard the fact that we could actually
					// do cross-compilation.
					// We have the target platform in `input.Input.Contract.Data.GolangSourceData.Platforms`
					Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
					Filename:  input.Input.Contract.Name,
					Version:   input.Input.Contract.Version,
					DependsOn: input.Input.Contract.Data.GolangSourceData.DependsOn,
				}},
			},
		}
		if err := build(inputArtifactPath, debug, BuildOpts{
			BinaryName: input.Input.Contract.Name,
			Version:    input.Input.Contract.Version,
			Tags:       input.Input.Contract.Data.GolangSourceData.Tags,
			OutputDir:  outputArtifactPath,
		}); err != nil {
			log.Fatalf("build failed: %s", err)
		}
		output.Results = append(output.Results, result)

	case "test":
		log.Fatalf("unimplemented mode %q", mode)

	default:
		log.Fatalf("unknown mode %q", mode)

	}

	// write input contract
	outputFd, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(outputFd).Encode(output); err != nil {
		log.Fatal(err)
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
