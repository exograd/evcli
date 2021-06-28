package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
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

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var apiErr Error
		if err := json.Unmarshal(resBody, &apiErr); err == nil {
			return &apiErr
		}

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
	query.Add("size", "100")
	uri := &url.URL{Path: "/v0/projects", RawQuery: query.Encode()}

	err := c.SendRequest("GET", uri, nil, &page)
	if err != nil {
		return nil, err
	}

	return page.Elements, nil
}

func (c *Client) FetchProjectByName(name string) (*Project, error) {
	uri := &url.URL{Path: "/v0/projects/name/" + url.QueryEscape(name)}

	var project Project

	err := c.SendRequest("GET", uri, nil, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) DeleteProject(id string) error {
	uri := &url.URL{Path: "/v0/projects/id/" + url.QueryEscape(id)}

	return c.SendRequest("DELETE", uri, nil, nil)
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
	res, err := rt.RoundTripper.RoundTrip(req)
	trace("%s %s", req.Method, req.URL.String())
	return res, err
}
