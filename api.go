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

	"github.com/pkg/errors"
)

type Error struct {
	ErrorString string      `json:"error"`
	Code        string      `json:"code"`
	Data        interface{} `json:"data,omitempty"`
}

func (err Error) Error() string {
	return err.ErrorString
}

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
		return nil, errors.Wrap(err, "invalid api endpoint")
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
			return errors.Wrap(err, "cannot encode body")
		}

		bodyReader = bytes.NewReader(bodyData)
	}

	req, err := http.NewRequest(method, uri.String(), bodyReader)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot send request")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "cannot read response body")
		}

		var apiErr Error
		if err := json.Unmarshal(resBody, &apiErr); err == nil {
			return &apiErr
		}

		return fmt.Errorf("request failed with status %d: %s",
			res.StatusCode, string(resBody))
	}

	return err
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
