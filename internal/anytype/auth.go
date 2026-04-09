package anytype

import (
	"context"
	"fmt"

	anytypeSdk "github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"
)

const appName = "google_keep_anytype_migration"

func AuthWithChallenge(ctx context.Context, addr, code string) (anytypeSdk.Client, error) {
	client := anytypeSdk.NewClient(anytypeSdk.WithBaseURL(addr))
	authClient := client.Auth()

	challengeResp, err := authClient.CreateChallenge(ctx, appName)
	if err != nil {
		return nil, err
	}

	if challengeResp == nil {
		return nil, fmt.Errorf("create challenge request returned empty response")
	}

	apiKeyResp, err := authClient.CreateApiKey(ctx, challengeResp.ChallengeID, code)
	if err != nil {
		return nil, err
	}

	if apiKeyResp == nil {
		return nil, fmt.Errorf("create api key request returned empty response")
	}

	return anytypeSdk.NewClient(anytypeSdk.WithBaseURL(addr), anytypeSdk.WithAppKey(apiKeyResp.ApiKey)), nil
}
