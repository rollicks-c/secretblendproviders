package prompt

import (
	"fmt"
	"github.com/rollicks-c/secretblend"
	"github.com/rollicks-c/term"
	"strings"
)

type Secret map[string]string

type Client struct {
}

type Option func(c *Client)

func Register() error {
	secretblend.AddProvider(NewClient(), "prompt://")
	return nil
}

func NewClient(options ...Option) *Client {
	client := &Client{}
	for _, opt := range options {
		opt(client)
	}
	return client
}

func (c Client) LoadSecret(path string) (string, error) {

	prompt := fmt.Sprintf("enter value for [%s]", path)
	value, err := term.PromptSecret(prompt)
	if err != nil {
		return "", err
	}

	value = strings.TrimSpace(value)
	return value, nil
}
