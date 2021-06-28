package main

type Error struct {
	ErrorString string      `json:"error"`
	Code        string      `json:"code"`
	Data        interface{} `json:"data,omitempty"`
}

func (err Error) Error() string {
	return err.ErrorString
}

type APIStatus struct {
}

type Cursor struct {
	Before  string `json:"before,omitempty"`
	After   string `json:"after,omitempty"`
	Size    uint   `json:"size,omitempty"`
	Reverse bool   `json:"reverse"`
}

type ProjectPage struct {
	Elements []*Project `json:"elements"`
	Previous *Cursor    `json:"previous,omitempty"`
	Next     *Cursor    `json:"next,omitempty"`
}

type Project struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
