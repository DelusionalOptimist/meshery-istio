package config

import (
	"fmt"
	"os"
	"path"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshkit/utils"
	"gopkg.in/yaml.v2"
)

var generatedTemplatePath = path.Join(RootPath(), "generated-imagehub-filter.yaml")

// Structs generated from template/imagehub-filter.yaml
// EnvoyFilter
type ImageHubTemplate struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// VmConfig
type VmConfig struct {
	Code             Code   `yaml:"code"`
	Runtime          string `yaml:"runtime"`
	VmId             string `yaml:"vmId"`
	AllowPrecompiled bool   `yaml:"allow_precompiled"`
}

// Labels
type Labels struct {
	App     string `yaml:"app"`
	Version string `yaml:"version"`
}

// Metadata
type Metadata struct {
	Name string `yaml:"name"`
}

// ConfigPatches
type ConfigPatches struct {
	ApplyTo string `yaml:"applyTo"`
	Match   Match  `yaml:"match"`
	Patch   Patch  `yaml:"patch"`
}

// Match
type Match struct {
	Context  string   `yaml:"context"`
	Proxy    Proxy    `yaml:"proxy"`
	Listener Listener `yaml:"listener"`
}

// Proxy
type Proxy struct {
	ProxyVersion string `yaml:"proxyVersion"`
}

// Value
type Value struct {
	Name        string      `yaml:"name"`
	TypedConfig TypedConfig `yaml:"typed_config"`
}

// Config
type FilterConfig struct {
	Configuration Configuration `yaml:"configuration"`
	RootId        string        `yaml:"root_id"`
	VmConfig      VmConfig      `yaml:"vmConfig"`
}

// Configuration
type Configuration struct {
	Type2 string `yaml:"@type"`
	Value string `yaml:"value"`
}

// Local
type Local struct {
	Filename string `yaml:"filename"`
}

// Listener
type Listener struct {
	PortNumber  int         `yaml:"portNumber"`
	FilterChain FilterChain `yaml:"filterChain"`
}

// FilterChain
type FilterChain struct {
	Filter Filter `yaml:"filter"`
}

// Filter
type Filter struct {
	Name      string    `yaml:"name"`
	SubFilter SubFilter `yaml:"subFilter"`
}

// TypedConfigValue
type TypedConfigValue struct {
	Config FilterConfig `yaml:"config"`
}

// Code
type Code struct {
	Local Local `yaml:"local"`
}

// WorkloadSelector
type WorkloadSelector struct {
	Labels Labels `yaml:"labels"`
}

// Spec
type Spec struct {
	ConfigPatches    []ConfigPatches  `yaml:"configPatches"`
	WorkloadSelector WorkloadSelector `yaml:"workloadSelector"`
}

// SubFilter
type SubFilter struct {
	Name string `yaml:"name"`
}

// Patch
type Patch struct {
	Operation string `yaml:"operation"`
	Value     Value  `yaml:"value"`
}

// TypedConfig
type TypedConfig struct {
	Type1   string           `yaml:"@type"`
	TypeUrl string           `yaml:"type_url"`
	Value   TypedConfigValue `yaml:"value"`
}

// function GenerateImageHubTemplate() generates a template with some changes
// and returns its path alongwith an error(if any) from the given file url
// (http, https, file).
func GenerateImagehubTemplate(t adapter.Template) (string, error) {

	// returns contents of the file in the url
	templateContent, err := utils.ReadFileSource(string(t))
	if err != nil {
		return "", err
	}

	ef := ImageHubTemplate{}

	// Unmarshalling the yaml string into an EnvoyFilter object
	err = yaml.Unmarshal([]byte(templateContent), &ef)
	if err != nil {
		return "", err
	}

	// change something
	ef.Spec.ConfigPatches[0].Patch.Value.TypedConfig.Value.Config.Configuration.Value = "some-filter"

	// Marshalling to yaml after making changes
	newYaml, err := yaml.Marshal(ef)
	if err != nil {
		return "", err
	}

	// checking if previously generated template exists. If yes then delete it.
	if _, err := os.Stat(generatedTemplatePath); !os.IsNotExist(err) {
		err = os.Remove(generatedTemplatePath)
		if err != nil {
			return "", err
		}
	}
	fmt.Println(os.Getwd())

	_ = utils.CreateFile(newYaml, "generated-imagehub-filter.yaml", RootPath())

	return generatedTemplatePath, nil
}
