package rest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	baseUrl = "http://localhost:31009"
)

func TestChallenge(t *testing.T) {
	client, err := NewClient(baseUrl)
	require.NoError(t, err)

	r, err := client.CreateChallenge(context.Background())
	require.NoError(t, err)
	require.NotNil(t, r)
	require.NotEmpty(t, r.ChallengeId)
}

func TestCreateApiKey(t *testing.T) {
	client, err := NewClient(baseUrl)
	require.NoError(t, err)

	ctx := context.Background()

	r, err := client.CreateChallenge(ctx)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.NotEmpty(t, r.ChallengeId)

	var code string
	_, err = fmt.Scanln(&code)
	require.NoError(t, err)

	apiResp, err := client.CreateApiKey(ctx, r.ChallengeId, code)
	require.NoError(t, err)
	require.NotNil(t, apiResp)
	require.NotEmpty(t, apiResp.ApiKey)
}
