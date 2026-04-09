package state

import "github.com/epheo/anytype-go"

type AppAuthState struct {
	APIAddress string
	APIKey     string
	client     anytype.Client
}

type AppAuthStater interface {
	GetAPIAddress() string
	SetAPIAddress(addr string)
	GetAPIKey() string
	SetAPIKey(key string)
	GetClient() anytype.Client
	SetClient(client anytype.Client)
}

func (as *AppAuthState) GetAPIAddress() string {
	return as.APIAddress
}

func (as *AppAuthState) SetAPIAddress(addr string) {
	as.APIAddress = addr
}

func (as *AppAuthState) GetAPIKey() string {
	return as.APIKey
}

func (as *AppAuthState) SetAPIKey(key string) {
	as.APIKey = key
}

func (as *AppAuthState) GetClient() anytype.Client {
	return as.client
}

func (as *AppAuthState) SetClient(client anytype.Client) {
	as.client = client
}
