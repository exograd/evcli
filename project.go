package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

type ProjectFile struct {
	Name string `json:"name"`
}

func LoadProjectFile(dirPath string) (*ProjectFile, error) {
	filePath := path.Join(dirPath, "eventline-project.json")

	trace("loading project file from %s", filePath)

	var pf ProjectFile
	if err := pf.LoadFile(filePath); err != nil {
		return nil, err
	}

	return &pf, nil
}

func (pf *ProjectFile) LoadFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	return pf.LoadData(data)
}

func (pf *ProjectFile) LoadData(data []byte) error {
	if err := json.Unmarshal(data, pf); err != nil {
		return fmt.Errorf("cannot parse json data: %w", err)
	}

	return nil
}

func (pf *ProjectFile) WriteFile(dirPath string) error {
	data, err := json.Marshal(pf)
	if err != nil {
		return fmt.Errorf("cannot encode data: %w", err)
	}

	filePath := path.Join(dirPath, "eventline-project.json")

	return ioutil.WriteFile(filePath, data, 0644)
}
