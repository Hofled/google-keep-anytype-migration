package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype"
)

type CreateObjectRequest struct {
	TypeKey    string                          `json:"type_key"`
	Body       string                          `json:"body,omitempty"`
	Icon       *anytype.Icon                   `json:"icon,omitempty"`
	Name       string                          `json:"name,omitempty"`
	Properties []anytype.PropertyLinkWithValue `json:"properties,omitempty"`
	TemplateID string                          `json:"template_id,omitempty"`
}

type CreatedObjectResponse struct {
	Object Object `json:"object"`
}

type Object struct {
	ID               string                          `json:"id"`
	Name             string                          `json:"name"`
	Icon             *anytype.Icon                   `json:"icon,omitempty"`
	Body             string                          `json:"body,omitempty"`
	TypeKey          string                          `json:"type_key"`
	Archived         bool                            `json:"archived"`
	CreatedDate      int64                           `json:"created_date"`
	LastModifiedDate int64                           `json:"last_modified_date"`
	LastOpenedDate   int64                           `json:"last_opened_date"`
	Properties       []anytype.PropertyLinkWithValue `json:"properties,omitempty"`
}

const createObjectEndpoint = "spaces/%s/objects"

func (c *Client) CreateObject(ctx context.Context, spaceId string, createObjReq CreateObjectRequest) (*CreatedObjectResponse, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&createObjReq); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(createObjectEndpoint, url.PathEscape(spaceId))

	req, err := c.newRequest(ctx, http.MethodPost, endpoint, buf)
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

	var createdObjResp CreatedObjectResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&createdObjResp); decodeErr != nil {
		return nil, decodeErr
	}

	return &createdObjResp, nil
}
