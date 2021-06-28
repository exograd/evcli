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
