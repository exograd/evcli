package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

type ProjectFile struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
}

func (pf *ProjectFile) Read(dirPath string) error {
	filePath := path.Join(dirPath, "eventline-project.json")

	trace("reading project file %s", filePath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	if err := json.Unmarshal(data, pf); err != nil {
		return fmt.Errorf("cannot decode json data: %w", err)
	}

	return nil
}

func (pf *ProjectFile) Write(dirPath string) error {
	data, err := json.Marshal(pf)
	if err != nil {
		return fmt.Errorf("cannot encode json data: %w", err)
	}

	filePath := path.Join(dirPath, "eventline-project.json")

	trace("writing project file %s", filePath)

	return ioutil.WriteFile(filePath, data, 0644)
}
