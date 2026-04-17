package rest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSpaces(t *testing.T) {
	data := `
	{
  "data": [
    {
      "description": "The local-first wiki",
      "gateway_url": "http://127.0.0.1:31006",
      "icon": {
        "emoji": "📄",
        "format": "icon"
      },
      "id": "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1",
      "name": "My Space",
      "network_id": "N83gJpVd9MuNRZAuJLZ7LiMntTThhPc6DtzWWVjb1M3PouVU",
      "object": "space"
    },
    {
      "description": "The local-first wiki",
      "gateway_url": "http://127.0.0.1:31006",
      "icon": {
        "file": "bafybeieptz5hvcy6txplcvphjbbh5yjc2zqhmihs3owkh5oab4ezauzqay",
        "format": "file"
      },
      "id": "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1",
      "name": "My Space",
      "network_id": "N83gJpVd9MuNRZAuJLZ7LiMntTThhPc6DtzWWVjb1M3PouVU",
      "object": "space"
    }
  ],
  "pagination": {
    "has_more": true,
    "limit": 100,
    "offset": 0,
    "total": 1000
  }
}`

	var spacesResp ListSpacesResponse
	err := json.Unmarshal([]byte(data), &spacesResp)
	require.NoError(t, err)
	require.NotEmpty(t, spacesResp.Data)
}
