package oam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	basePath, _  = os.Getwd()
	workloadPath = filepath.Join(basePath, "oam", "workloads")
	traitPath    = filepath.Join(basePath, "oam", "traits")
)

type OAMGenericStructure struct {
	OAMDefinition interface{} `json:"oam_definition,omitempty"`
	OAMRefSchema  string      `json:"oam_ref_schema,omitempty"`
	Host          string      `json:"host,omitempty"`
}

// RegisterWorkloads will register all of the workload definitions
// present in the path oam/workloads
//
// Registeration process will send POST request to $runtime/api/experimental/oam/workload
func RegisterWorkloads(runtime, host string) error {
	workloads := []string{
		"istiomesh",
		"grafanaistioaddon",
		"prometheusistioaddon",
		"zipkinistioaddon",
		"jaegeristioaddon",
	}

	for _, workload := range workloads {
		oamcontent, err := readDefintionAndSchema(workloadPath, workload)
		if err != nil {
			return err
		}
		oamcontent.Host = host

		// Convert struct to json
		byt, err := json.Marshal(oamcontent)
		if err != nil {
			return err
		}

		reader := bytes.NewReader(byt)

		if err := register(fmt.Sprintf("%s/api/experimental/oam/%s", runtime, "workload"), reader); err != nil {
			return err
		}
	}

	return nil
}

// RegisterTraits will register all of the trait definitions
// present in the path oam/traits
//
// Registeration process will send POST request to $runtime/api/experimental/oam/trait
func RegisterTraits(runtime, host string) error {
	traits := []string{
		"automaticsidecarinjection",
		"mtls",
	}

	for _, trait := range traits {
		oamcontent, err := readDefintionAndSchema(traitPath, trait)
		if err != nil {
			return err
		}
		oamcontent.Host = host

		// Convert struct to json
		byt, err := json.Marshal(oamcontent)
		if err != nil {
			return err
		}

		reader := bytes.NewReader(byt)

		if err := register(fmt.Sprintf("%s/api/experimental/oam/%s", runtime, "trait"), reader); err != nil {
			return err
		}
	}

	return nil
}

func readDefintionAndSchema(path, name string) (*OAMGenericStructure, error) {
	definitionName := fmt.Sprintf("%s_definition.json", name)
	schemaName := fmt.Sprintf("%s.meshery.layer5.io.schema.json", name)

	definition, err := ioutil.ReadFile(filepath.Join(path, definitionName))
	if err != nil {
		return nil, err
	}

	var definitionMap map[string]interface{}
	if err := json.Unmarshal(definition, &definitionMap); err != nil {
		return nil, err
	}

	schema, err := ioutil.ReadFile(filepath.Join(path, schemaName))
	if err != nil {
		return nil, err
	}

	return &OAMGenericStructure{OAMDefinition: definitionMap, OAMRefSchema: string(schema)}, nil
}

func register(host string, content io.Reader) error {
	resp, err := http.Post(host, "application/json", content)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"register process failed, host returned status: %s with status code %d",
			resp.Status,
			resp.StatusCode,
		)
	}

	return nil
}
