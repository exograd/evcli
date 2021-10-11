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

func (pf *ProjectFile) UnmarshalJSON(data []byte) error {
	type ProjectFile2 ProjectFile

	pf2 := ProjectFile2(*pf)
	if err := json.Unmarshal(data, &pf2); err != nil {
		return err
	}

	if pf2.Id == "" {
		return fmt.Errorf("missing project id")
	}

	*pf = ProjectFile(pf2)
	return nil
}

func (pf *ProjectFile) Read(dirPath string) error {
	filePath := path.Join(dirPath, "eventline-project.json")

	p.Debug(1, "reading project file %s", filePath)
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

	p.Debug(1, "writing project file %s", filePath)

	return ioutil.WriteFile(filePath, data, 0644)
}
