package envvar

import (
	"github.com/rollicks-c/secretblend"
)

type Client struct {
	ignoreMissing bool
}

type Option func(c *Client)

func Register(options ...Option) error {

	// register client
	op := NewClient(options...)
	secretblend.AddProvider(op, "envvar://")

	return nil
}

func RegisterGlobally(options ...Option) error {

	// register client
	cl := NewClient(options...)
	secretblend.AddGlobalProvider(cl)

	return nil
}

func WithIgnoreMissing(ignoreMissing bool) Option {
	return func(c *Client) {
		c.ignoreMissing = ignoreMissing
	}
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

	value, err := c.injectVars(uri)
	if err != nil {
		return "", err
	}

	return value, nil
}
