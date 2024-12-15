package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var (
	//go:embed api.schema.yaml
	apiSchemaYAML []byte
)

func validateAPIFile(apiFile string) error {

	// Load the JSON Schema
	apiSchemaJSONFile, err := createAPISchemaJSONFile()
	if err != nil {
		return fmt.Errorf("failed to create API schema JSON file: %w", err)
	}
	defer os.Remove(apiSchemaJSONFile)

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + apiSchemaJSONFile)

	// Open the API YAML file
	f, err := os.Open(apiFile)
	if err != nil {
		return fmt.Errorf("failed to open API file '%s': %w", apiFile, err)
	}
	defer f.Close()

	// Convert YAML content to JSON
	var jsonAPIData bytes.Buffer
	if err := yamlToJSON(f, &jsonAPIData); err != nil {
		return fmt.Errorf("failed to convert API YAML to JSON: %w", err)
	}

	// Validate the API JSON data against the schema
	documentLoader := gojsonschema.NewBytesLoader(jsonAPIData.Bytes())
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed to validate API file against schema: %w", err)
	}

	// Check validation result
	if !result.Valid() {
		validationErrors := []error{fmt.Errorf("API file '%s' is invalid", apiFile)}
		for _, err := range result.Errors() {
			validationErrors = append(validationErrors, errors.New(err.String()))
		}
		return errors.Join(validationErrors...)
	}

	// File is valid
	return nil
}

// createAPISchemaJSONFile generates a temporary JSON file containing the API schema converted from an embedded YAML file.
// It returns the file path to the created JSON file or an error if the creation or conversion process fails.
func createAPISchemaJSONFile() (string, error) {
	apiSchemaJSONFile, err := os.CreateTemp("", "api.schema.*.json")
	if err != nil {
		return "", fmt.Errorf("failed creating JSON file: %w", err)
	}
	defer apiSchemaJSONFile.Close()

	if err := yamlToJSON(bytes.NewReader(apiSchemaYAML), apiSchemaJSONFile); err != nil {
		return "", fmt.Errorf("failed converting API schema (in YAML) into a JSON Schema: %w", err)
	}

	return apiSchemaJSONFile.Name(), nil
}

// yamlToJSON converts YAML content from the source reader to JSON format and writes it to the target writer.
// It ensures compatibility by transforming YAML maps with non-string keys into JSON-compatible maps.
// Returns an error if decoding the YAML or encoding the JSON fails.
func yamlToJSON(source io.Reader, target io.Writer) error {
	var yamlData interface{}

	yamlDecoder := yaml.NewDecoder(source)
	if err := yamlDecoder.Decode(&yamlData); err != nil {
		return fmt.Errorf("failed decoding YAML file '%s': %w", source, err)
	} else {
		// Important: Convert map[interface{}]interface{} to map[string]interface{} for JSON compatibility
		yamlData = convertKeysToStrings(yamlData)
	}

	//f, err := os.OpenFile(jsonFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	//if err != nil {
	//	return fmt.Errorf("failed to open JSON file '%s' for writing: %w", jsonFile, err)
	//}
	//defer f.Close()

	jsonEncoder := json.NewEncoder(target)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(yamlData); err != nil {
		return fmt.Errorf("failed encoding YAML data from '%s' into JSON: %w", source, err)
	}

	return nil
}

// convertKeysToStrings converts map keys from any type to strings recursively for JSON compatibility.
// If the input is a slice, it processes each element recursively. Returns the modified structure.
func convertKeysToStrings(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		mapped := map[string]interface{}{}
		for key, val := range v {
			mapped[fmt.Sprintf("%v", key)] = convertKeysToStrings(val)
		}
		return mapped
	case []interface{}:
		for i, val := range v {
			v[i] = convertKeysToStrings(val)
		}
	}
	return value
}
