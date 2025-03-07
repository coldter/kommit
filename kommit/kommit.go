package kommit

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseUrl *url.URL
	token   string

	// Single instance to be reused by all clients
	base *client

	Ai *AiClient
}

// Client struct that will be aliases by all other clients
type client struct {
	client *Client
}

// APIResponse represents the common response structure
type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

func New(base *url.URL, token string) *Client {
	c := &Client{baseUrl: base, token: token}

	c.base = &client{c}
	c.Ai = (*AiClient)(c.base)
	return c
}

func (t *Client) newRequest(method, urlPath string, body io.Reader, headers ...map[string]string) (*http.Request, error) {
	u, err := url.Parse(t.baseUrl.String())
	if err != nil {
		return nil, err
	}
	u, err = u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if t.token != "" {
		req.Header.Add("Authorization", fmt.Sprint("Bearer ", t.token))
	}

	if len(headers) > 0 {
		for _, h := range headers {
			for k, v := range h {
				req.Header.Add(k, v)
			}
		}
	}

	return req, nil
}

func (t *Client) do(method, path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	req, err := t.newRequest(method, path, body, headers...)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *Client) Get(path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return t.do("GET", path, body, headers...)
}

func (t *Client) Post(path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return t.do("POST", path, body, headers...)
}

func (t *Client) Patch(path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return t.do("PATCH", path, body, headers...)
}

func (t *Client) Put(path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return t.do("PUT", path, body, headers...)
}

func (t *Client) Delete(path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return t.do("DELETE", path, body, headers...)
}
