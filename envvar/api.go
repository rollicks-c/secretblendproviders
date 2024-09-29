package envvar

import (
	"fmt"
	"github.com/rollicks-c/secretblend"
	"os"
)

type Client struct {
}

type Option func(c *Client)

func Register() error {

	// register client
	op := NewClient()
	secretblend.AddProvider(op, "envvar://")

	return nil
}

func NewClient(options ...Option) *Client {

	// create lient
	client := &Client{}

	// apply options
	for _, opt := range options {
		opt(client)
	}

	return client
}

func (c Client) LoadSecret(uri string) (string, error) {

	value, ok := os.LookupEnv(uri)
	if !ok {
		return "", fmt.Errorf("failed to load secret - envvar not provided: %s", uri)
	}

	return value, nil
}
