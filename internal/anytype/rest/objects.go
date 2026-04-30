package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype"
)

type CreateObjectRequest struct {
	Name       string                          `json:"name"`
	Icon       *anytype.Icon                   `json:"icon,omitempty"`
	Body       string                          `json:"body"`
	TemplateId string                          `json:"template_id"`
	TypeKey    string                          `json:"type_key"`
	Properties []anytype.PropertyLinkWithValue `json:"properties,omitempty"`
}

type CreatedObjectResponse struct {
	Object Object `json:"object"`
}

type Object struct {
	Object     string       `json:"object"`
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Icon       anytype.Icon `json:"icon"`
	Archived   bool         `json:"archived"`
	SpaceId    string       `json:"space_id"`
	Snippet    string       `json:"snippet"`
	Layout     string       `json:"layout"`
	Type       *Type        `json:"type"`
	Properties []Property   `json:"properties"`
	Markdown   string       `json:"markdown"`
}

type Type struct {
	Object     string        `json:"object"`
	Id         string        `json:"id"`
	Key        string        `json:"key"`
	Name       string        `json:"name"`
	PluralName string        `json:"plural_name"`
	Icon       *anytype.Icon `json:"icon,omitempty"`
	Archived   bool          `json:"archived"`
	Layout     string        `json:"layout"`
	Properties []Property    `json:"properties"`
}

type Tag struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Property struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Format      string   `json:"format"`
	Text        *string  `json:"text,omitempty"`
	Number      *float64 `json:"number,omitempty"`
	Select      *Tag     `json:"select,omitempty"`
	MultiSelect []Tag    `json:"multi_select,omitempty"`
	Date        *string  `json:"date,omitempty"`
	File        []string `json:"file,omitempty"`
	Checkbox    *bool    `json:"checkbox,omitempty"`
	Url         *string  `json:"url,omitempty"`
	Email       *string  `json:"email,omitempty"`
	Phone       *string  `json:"phone,omitempty"`
	Object      string   `json:"object,omitempty"`
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
		body, _ := io.ReadAll(resp.Body)
		return nil, &InvalidResponseErr{statusCode: resp.StatusCode, body: string(body)}
	}

	var createdObjResp CreatedObjectResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&createdObjResp); decodeErr != nil {
		return nil, decodeErr
	}

	return &createdObjResp, nil
}
