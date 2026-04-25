package rest

import (
	"context"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
	apiVersionDate = "2025-11-08"
	apiVersion     = "v1"
)

type Client struct {
	baseUrl    *url.URL
	apiKey     string
	httpClient http.Client
}

func NewClient(baseUrl string) (*Client, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseUrl: parsedUrl,
		httpClient: http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

func (c *Client) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *Client) newRequest(ctx context.Context, method, apiPath string, bodyReader io.Reader) (*http.Request, error) {
	parsedPath, err := url.Parse(apiPath)
	if err != nil {
		return nil, err
	}

	reqUrl := c.baseUrl.JoinPath(apiVersion, parsedPath.Path)
	reqUrl.RawQuery = parsedPath.RawQuery

	req, err := http.NewRequestWithContext(ctx, method, reqUrl.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Anytype-Version", apiVersionDate)

	if len(c.apiKey) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	return req, nil
}

type InvalidResponseErr struct {
	Resp *http.Response
}

func (ire InvalidResponseErr) Error() string {
	return fmt.Sprintf("invalid response: %s", ire.Resp.Status)
}
