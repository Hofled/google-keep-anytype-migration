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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, &InvalidResponseErr{resp}
	}

	var challengeResp ChallengeResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&challengeResp); decodeErr != nil {
		return nil, decodeErr
	}

	return &challengeResp, nil
}

type CreateApiKeyRequest struct {
	ChallengeId string `json:"challenge_id"`
	Code        string `json:"code"`
}

type CreateApiKeyResponse struct {
	ApiKey string `json:"api_key"`
}

const createApiKeyEndpoint = "auth/api_keys"

func (c *Client) CreateApiKey(ctx context.Context, challengeId, code string) (*CreateApiKeyResponse, error) {
	apiKeyReq := CreateApiKeyRequest{
		ChallengeId: challengeId,
		Code:        code,
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&apiKeyReq); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, createApiKeyEndpoint, buf)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, &InvalidResponseErr{resp}
	}

	var apiKeyResp CreateApiKeyResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&apiKeyResp); decodeErr != nil {
		return nil, decodeErr
	}

	return &apiKeyResp, nil
}
