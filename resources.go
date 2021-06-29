package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ResourceSet struct {
	Resources []*Resource   `json:"-"`
	Specs     []interface{} `json:"specs"`
}

type Resource struct {
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
		trace("loading %s", filePath)

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

			rs.Resources = append(rs.Resources, fileResource)
			rs.Specs = append(rs.Specs, jsonSpec)
		}
	}

	return nil
}

func LoadResourceFile(filePath string) ([]*Resource, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	resources := []*Resource{}
	document := 1

	for {
		var value interface{}
		if err := decoder.Decode(&value); err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("cannot decode yaml data: %w", err)
		}

		resource := Resource{
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

			if match == false {
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
