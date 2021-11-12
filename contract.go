package main

const (
	TypeGolangSource = "type-product-os-t-golang-source@1.0.1"
	TypeExecutable   = "type-product-os-t-executable@1.0.1"
)

type NInput struct {
	Input TransformerAsset `json:"input"`
}

type NOutput struct {
	Results []TransformerAsset `json:"results"`
}

type TransformerAsset struct {
	Contract     Contract `json:"contract"`
	ArtifactPath string   `json:"artifactPath"`
}

type Contract struct {
	Type    string       `json:"type"`
	Version string       `json:"version,omitempty"`
	Name    string       `json:"name,omitempty"`
	Data    ContractData `json:"data"`
}

type ContractData struct {
	GolangSourceData
	ExecutableData
}

type GolangSourceData struct {
	Platforms []string            `json:"platforms,omitempty"`
	Tags      []string            `json:"tags,omitempty"`
	DependsOn map[string][]string `json:"dependsOn,omitempty"`
}

type ExecutableData struct {
	Platform  string              `json:"platform,omitempty"`
	Filename  string              `json:"filename,omitempty"`
	Version   string              `json:"version,omitempty"`
	DependsOn map[string][]string `json:"dependsOn,omitempty"`
}
