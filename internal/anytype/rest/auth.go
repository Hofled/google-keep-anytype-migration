package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/consts"
)

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
