package main

import (
	"encoding/json"
	"fmt"
	"time"
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

type ProjectSearchQuery struct {
	Id []string `json:"id"`
}

type Project struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Projects []*Project

func (ps Projects) GroupById() map[string]*Project {
	table := make(map[string]*Project)
	for _, p := range ps {
		table[p.Id] = p
	}

	return table
}

type PipelinePage struct {
	Elements []*Pipeline `json:"elements"`
	Previous *Cursor     `json:"previous,omitempty"`
	Next     *Cursor     `json:"next,omitempty"`
}

type Pipeline struct {
	Id           string     `json:"id,omitempty"`
	Name         string     `json:"name"`
	OrgId        string     `json:"org_id"`
	ProjectId    string     `json:"project_id,omitempty"`
	CreationTime time.Time  `json:"creation_time"`
	PipelineId   string     `json:"pipeline_id,omitempty"`
	TriggerId    string     `json:"trigger_id,omitempty"`
	EventId      string     `json:"event_id,omitempty"`
	Concurrent   bool       `json:"concurrent,omitempty"`
	Status       string     `json:"status"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
}

func (p *Pipeline) Duration() *time.Duration {
	if p.StartTime == nil || p.EndTime == nil {
		return nil
	}

	d := p.EndTime.Sub(*p.StartTime)
	return &d
}

type Pipelines []*Pipeline

func (ps Pipelines) ProjectIds() []string {
	idTable := make(map[string]struct{})
	for _, p := range ps {
		idTable[p.ProjectId] = struct{}{}
	}

	ids := make([]string, 0, len(idTable))
	for id := range idTable {
		ids = append(ids, id)
	}

	return ids
}

type TaskPage struct {
	Elements []*Task `json:"elements"`
	Previous *Cursor `json:"previous,omitempty"`
	Next     *Cursor `json:"next,omitempty"`
}

type Task struct {
	Id             string     `json:"id,omitempty"`
	OrgId          string     `json:"org_id"`
	ProjectId      string     `json:"project_id,omitempty"`
	PipelineId     string     `json:"pipeline_id,omitempty"`
	TaskId         string     `json:"task_id,omitempty"`
	InstanceId     int        `json:"instance_id,omitempty"`
	Status         string     `json:"status"`
	StartTime      *time.Time `json:"start_time,omitempty"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	FailureMessage string     `json:"failure_message,omitempty"`
}

type Tasks []*Task
