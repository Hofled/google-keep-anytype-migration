package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/consts"
)

const (
	defaultTimeout = 5 * time.Second
	apiVersionDate = "2025-11-08"
	apiVersion     = "v1"
)

type Client struct {
	baseUrl *url.URL
	apiKey  string
}

func NewClient(baseUrl string) (*Client, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseUrl: parsedUrl,
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl.String(), bodyReader)
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

type ChallengeRequest struct {
	AppName string `json:"app_name"`
}

type ChallengeResponse struct {
	ChallengeId string `json:"challenge_id"`
}

const challengeEndpoint = "auth/challenges"

func (c *Client) CreateChallenge(ctx context.Context) (*ChallengeResponse, error) {
	challengeReq := ChallengeRequest{
		AppName: consts.AppName,
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&challengeReq); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, challengeEndpoint, buf)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var challengeResp ChallengeResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&challengeResp); decodeErr != nil {
		return nil, decodeErr
	}

	return &challengeResp, nil
}
