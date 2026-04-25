package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype"
)

type Space struct {
	Description string       `json:"description"`
	GatewayUrl  string       `json:"gateway_url"`
	Icon        anytype.Icon `json:"icon"`
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	NetworkId   string       `json:"network_id"`
	Object      string       `json:"object"`
}

type ListSpacesResponse struct {
	Data []*Space `json:"data"`
}

const listSpacesEndpoint = "spaces"

func (c *Client) ListSpaces(ctx context.Context) (*ListSpacesResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, listSpacesEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &InvalidResponseErr{resp}
	}

	var listSpacesResp ListSpacesResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&listSpacesResp); decodeErr != nil {
		return nil, decodeErr
	}

	return &listSpacesResp, nil
}
