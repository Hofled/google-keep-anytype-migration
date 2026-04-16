package rest

import (
	"context"
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
