package main

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	ErrorString string          `json:"error"`
	Code        string          `json:"code"`
	RawData     json.RawMessage `json:"data,omitempty"`
	Data        interface{}     `json:"-"`
}

type InvalidRequestBodyError struct {
	JSVErrors []JSVError `json:"jsv_errors"`
}

type JSVError struct {
	Pointer string `json:"pointer"`
	Reason  string `json:"reason"`
}

func (err APIError) Error() string {
	return err.ErrorString
}

func (err *APIError) UnmarshalJSON(data []byte) error {
	type APIError2 APIError

	err2 := APIError2(*err)
	if err := json.Unmarshal(data, &err2); err != nil {
		return err
	}

	switch err2.Code {
	case "invalid_request_body":
		var errData InvalidRequestBodyError

		if err := json.Unmarshal(err2.RawData, &errData); err != nil {
			return fmt.Errorf("invalid jsv errors: %w", err)
		}

		err2.Data = errData
	}

	*err = APIError(err2)
	return nil
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
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
