package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func formatAddObjectsToListEndpoint(spaceId, listId string) string {
	return fmt.Sprintf("spaces/%s/lists/%s/objects", url.PathEscape(spaceId), url.PathEscape(listId))
}

type AddObjectsToListRequest struct {
	Objects []string `json:"objects"`
}

func (c *Client) AddObjectsToList(ctx context.Context, spaceId, listId string, objectIds []string) error {
	reqObj := AddObjectsToListRequest{objectIds}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&reqObj); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, formatAddObjectsToListEndpoint(spaceId, listId), buf)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &InvalidResponseErr{statusCode: resp.StatusCode, body: string(body)}
	}

	return nil
}
