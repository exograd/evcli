package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	APIKey    string
	ProjectId string

	httpClient *http.Client

	baseURI *url.URL
}

func NewClient(config *Config) (*Client, error) {
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: NewRoundTripper(http.DefaultTransport),
	}

	baseURI, err := url.Parse(config.API.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid api endpoint: %w", err)
	}

	client := &Client{
		httpClient: httpClient,

		baseURI: baseURI,
	}

	return client, nil
}

func (c *Client) SendRequest(method string, relURI *url.URL, body, dest interface{}) error {
	uri := c.baseURI.ResolveReference(relURI)

	var bodyReader io.Reader
	if body == nil {
		bodyReader = nil
	} else {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("cannot encode body: %w", err)
		}

		bodyReader = bytes.NewReader(bodyData)
	}

	req, err := http.NewRequest(method, uri.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	if c.ProjectId != "" {
		req.Header.Set("X-Eventline-Project-Id", c.ProjectId)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var apiErr APIError

		err := json.Unmarshal(resBody, &apiErr)
		if err == nil {
			return &apiErr
		}

		trace("cannot decode response body: %v", err)

		return fmt.Errorf("request failed with status %d: %s",
			res.StatusCode, string(resBody))
	}

	if dest != nil {
		if len(resBody) == 0 {
			return fmt.Errorf("empty response body")
		}

		if err := json.Unmarshal(resBody, dest); err != nil {
			return fmt.Errorf("cannot decode response body: %w", err)
		}
	}

	return err
}

func (c *Client) FetchProjects() ([]*Project, error) {
	var page ProjectPage

	query := url.Values{}
	query.Add("size", "20")
	uri := url.URL{Path: "/v0/projects", RawQuery: query.Encode()}

	err := c.SendRequest("GET", &uri, nil, &page)
	if err != nil {
		return nil, err
	}

	return page.Elements, nil
}

func (c *Client) FetchProjectByName(name string) (*Project, error) {
	uri := url.URL{Path: "/v0/projects/name/" + name}

	var project Project

	err := c.SendRequest("GET", &uri, nil, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) SearchProjects(searchQuery ProjectSearchQuery) (Projects, error) {
	var page ProjectPage

	query := url.Values{}
	query.Add("size", "20")
	uri := url.URL{Path: "/v0/projects/search", RawQuery: query.Encode()}

	err := c.SendRequest("POST", &uri, &searchQuery, &page)
	if err != nil {
		return nil, err
	}

	return page.Elements, nil
}

func (c *Client) CreateProject(project *Project) error {
	uri := url.URL{Path: "/v0/projects"}

	return c.SendRequest("POST", &uri, project, project)
}

func (c *Client) DeleteProject(id string) error {
	uri := url.URL{Path: "/v0/projects/id/" + id}

	return c.SendRequest("DELETE", &uri, nil, nil)
}

func (c *Client) DeployProject(id string, rs *ResourceSet) error {
	uri := url.URL{
		Path: "/v0/projects/id/" + id + "/resources",
	}

	return c.SendRequest("PUT", &uri, rs, nil)
}

func (c *Client) FetchCommands() ([]*Resource, error) {
	var commands Resources

	cursor := Cursor{Size: 20}

	for {
		var page ResourcePage

		query := url.Values{}

		query.Add("type", "command")
		query.Add("size", strconv.FormatUint(uint64(cursor.Size), 10))
		if cursor.After != "" {
			query.Add("after", cursor.After)
		}

		uri := url.URL{Path: "/v0/resources", RawQuery: query.Encode()}

		err := c.SendRequest("GET", &uri, nil, &page)
		if err != nil {
			return nil, err
		}

		commands = append(commands, page.Elements...)

		if page.Next == nil {
			break
		}

		cursor = *page.Next
	}

	return commands, nil
}

func (c *Client) FetchPipelines() (Pipelines, error) {
	var page PipelinePage

	query := url.Values{}
	query.Add("size", "20")
	query.Add("reverse", "")
	uri := url.URL{Path: "/v0/pipelines", RawQuery: query.Encode()}

	err := c.SendRequest("GET", &uri, nil, &page)
	if err != nil {
		return nil, err
	}

	return page.Elements, nil
}

func (c *Client) AbortPipeline(id string) error {
	uri := url.URL{Path: "/v0/pipelines/id/" + id + "/abort"}

	return c.SendRequest("POST", &uri, nil, nil)
}

func (c *Client) RestartPipeline(id string) error {
	uri := url.URL{Path: "/v0/pipelines/id/" + id + "/restart"}

	return c.SendRequest("POST", &uri, nil, nil)
}

func (c *Client) FetchTasks() (Tasks, error) {
	var page TaskPage

	query := url.Values{}
	query.Add("size", "20")
	query.Add("reverse", "")
	uri := url.URL{Path: "/v0/tasks", RawQuery: query.Encode()}

	err := c.SendRequest("GET", &uri, nil, &page)
	if err != nil {
		return nil, err
	}

	return page.Elements, nil
}

func (c *Client) CreateEvent(newEvent *NewEvent) (Events, error) {
	var events Events

	uri := url.URL{Path: "/v0/events"}
	err := c.SendRequest("POST", &uri, newEvent, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

type RoundTripper struct {
	http.RoundTripper
}

func NewRoundTripper(rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{
		RoundTripper: rt,
	}
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	res, err := rt.RoundTripper.RoundTrip(req)
	d := time.Now().Sub(start)

	var statusString string
	if res == nil {
		statusString = "-"
	} else {
		statusString = strconv.Itoa(res.StatusCode)
	}

	trace("%s %s %s %s", req.Method, req.URL.String(), statusString,
		FormatRequestDuration(d))
	return res, err
}

func FormatRequestDuration(d time.Duration) string {
	s := d.Seconds()

	switch {
	case s < 0.001:
		return fmt.Sprintf("%dμs", d.Microseconds())

	case s < 1.0:
		return fmt.Sprintf("%dms", d.Milliseconds())

	default:
		return fmt.Sprintf("%.1fs", s)
	}
}
