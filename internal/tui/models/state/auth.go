package state

import (
	"sync"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
)

type AppAuthState struct {
	apiAddress string
	apiKey     string
	client     *rest.Client

	mut sync.Mutex
}

type AppAuthStater interface {
	GetAPIAddress() string
	SetAPIAddress(addr string)
	GetAPIKey() string
	SetAPIKey(key string)
	GetClient() *rest.Client
}

func (as *AppAuthState) GetAPIAddress() string {
	return as.apiAddress
}

func (as *AppAuthState) SetAPIAddress(addr string) {
	as.apiAddress = addr
}

func (as *AppAuthState) GetAPIKey() string {
	return as.apiKey
}

func (as *AppAuthState) SetAPIKey(key string) {
	as.apiKey = key
}

func (as *AppAuthState) GetClient() *rest.Client {
	as.mut.Lock()
	defer as.mut.Unlock()

	if as.client != nil {
		return as.client
	}

	if len(as.apiAddress) == 0 || len(as.apiKey) == 0 {
		return nil
	}

	client, err := rest.NewClient(as.apiAddress)
	if err != nil {
		return nil
	}

	client.SetApiKey(as.apiKey)
	as.client = client

	return as.client
}
