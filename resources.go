package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/qri-io/jsonpointer"
	"gopkg.in/yaml.v3"
)

type ResourceSet struct {
	Resources []*ResourceFile `json:"-"`
	Specs     []interface{}   `json:"specs"`
}

type ResourceFile struct {
	Path     string
	Document int
	Value    interface{}
}

func (rs *ResourceSet) Load(dirPath string) error {
	extensions := []string{".yml", ".yaml"}
	filePaths, err := FindFiles(dirPath, extensions)
	if err != nil {
		return fmt.Errorf("cannot find files: %w", err)
	}

	for _, filePath := range filePaths {
		p.Debug(1, "loading resource file %s", filePath)

		fileResources, err := LoadResourceFile(filePath)
		if err != nil {
			return fmt.Errorf("cannot load %s: %w", filePath, err)
		}

		for _, fileResource := range fileResources {
			yamlSpec := fileResource.Value
			jsonSpec, err := YAMLValueToJSONValue(yamlSpec)
			if err != nil {
				return fmt.Errorf("%s: document %d is not a valid json "+
					"value: %w", filePath, fileResource.Document, err)
			}

			if SpecType(jsonSpec) == "task" {
				if err := LoadTaskSource(jsonSpec, dirPath); err != nil {
					return fmt.Errorf("%s: cannot load task source for "+
						"document %d: %w",
						fileResource.Path, fileResource.Document, err)
				}
			}

			rs.Resources = append(rs.Resources, fileResource)
			rs.Specs = append(rs.Specs, jsonSpec)
		}
	}

	return nil
}

func LoadResourceFile(filePath string) ([]*ResourceFile, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	resources := []*ResourceFile{}
	document := 1

	for {
		var value interface{}
		if err := decoder.Decode(&value); err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("cannot decode yaml data: %w", err)
		}

		resource := ResourceFile{
			Path:     filePath,
			Document: document,
			Value:    value,
		}

		resources = append(resources, &resource)
		document++
	}

	return resources, nil
}

func FindFiles(dirPath string, extensions []string) ([]string, error) {
	var filePaths []string

	err := filepath.Walk(dirPath,
		func(filePath string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ext := strings.ToLower(filepath.Ext(filePath))

			match := false
			for _, e := range extensions {
				if ext == e {
					match = true
					break
				}
			}

			if !match {
				return nil
			}

			filePaths = append(filePaths, filePath)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return filePaths, nil
}

func SpecType(spec interface{}) string {
	ptr, _ := jsonpointer.Parse("/type")

	value, err := ptr.Eval(spec)
	if err != nil {
		return ""
	}

	s, ok := value.(string)
	if !ok {
		return ""
	}

	return s
}

func LoadTaskSource(spec interface{}, dirPath string) error {
	ptr, _ := jsonpointer.Parse("/data/steps")

	values, err := ptr.Eval(spec)
	if err != nil {
		return nil
	}

	steps, ok := values.([]interface{})
	if !ok {
		return nil
	}

	for _, step := range steps {
		if err := LoadStepSource(step, dirPath); err != nil {
			return err
		}
	}

	return nil
}

func LoadStepSource(stepValue interface{}, dirPath string) error {
	step, ok := stepValue.(map[string]interface{})
	if !ok {
		return nil
	}

	sourceValue, found := step["source"]
	if !found {
		return nil
	}

	source, ok := sourceValue.(string)
	if !ok {
		return fmt.Errorf("%v is not a string", sourceValue)
	}

	sourcePath := path.Join(dirPath, source)

	p.Debug(1, "loading task step source file %s", sourcePath)

	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %w", sourcePath, err)
	}

	step["code"] = string(data)

	return nil
}

func (rf *ResourceFile) TypeAndName() (typeName string, name string) {
	value, ok := rf.Value.(map[string]interface{})
	if !ok {
		return "", ""
	}

	if sv, found := value["type"]; found {
		if s, ok := sv.(string); ok {
			typeName = s
		}
	}

	if sv, found := value["name"]; found {
		if s, ok := sv.(string); ok {
			name = s
		}
	}

	return
}
