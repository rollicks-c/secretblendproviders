package envvar

import (
	"github.com/rollicks-c/secretblend"
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

func RegisterGlobally() error {

	// register client
	cl := NewClient()
	secretblend.AddGlobalProvider(cl)

	return nil
}

func NewClient(options ...Option) *Client {

	// create client
	client := &Client{}

	// apply options
	for _, opt := range options {
		opt(client)
	}

	return client
}

func (c Client) LoadSecret(uri string) (string, error) {

	value, err := injectVars(uri)
	if err != nil {
		return "", err
	}

	return value, nil
}
