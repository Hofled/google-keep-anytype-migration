package state

import (
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
)

type AppAuthState struct {
	APIAddress string
	APIKey     string
	client     *rest.Client
}

type AppAuthStater interface {
	GetAPIAddress() string
	SetAPIAddress(addr string)
	GetAPIKey() string
	SetAPIKey(key string)
	GetClient() *rest.Client
	SetClient(client *rest.Client)
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

func (as *AppAuthState) GetClient() *rest.Client {
	return as.client
}

func (as *AppAuthState) SetClient(client *rest.Client) {
	as.client = client
}
