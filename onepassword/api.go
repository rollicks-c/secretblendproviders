package onepassword

import (
	"bytes"
	"fmt"
	"github.com/rollicks-c/secretblend"
	"os/exec"
	"strings"
)

type Client struct {
	cliPath string
}

type Option func(c *Client)

func Register() error {

	// use defaults
	cliPath := "op"

	// register client
	op := NewClient(cliPath)
	secretblend.AddProvider(op, "1password://")

	return nil
}

func NewClient(cliPath string, options ...Option) *Client {

	// create cli client
	client := &Client{
		cliPath: cliPath,
	}

	// apply options
	for _, opt := range options {
		opt(client)
	}

	return client
}

func (c Client) LoadSecret(uri string) (string, error) {

	// gather data
	parts := strings.Split(uri, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid uri: %s", uri)
	}
	itemID := parts[0]

	// sanity check
	field := parts[1]
	switch field {
	case "username":
	case "password":
	default:
		return "", fmt.Errorf("invalid 1password field: %s", field)
	}

	secret, err := c.GetItem(itemID, field)
	if err != nil {
		return "", err
	}

	return secret, nil
}

func (c Client) GetItem(id, field string) (string, error) {

	// build command
	cmd := exec.Command("op", "item", "get", id, "--field", field, "--reveal")
	var out bytes.Buffer
	cmd.Stdout = &out

	// run command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	// extract
	secret := out.String()
	secret = strings.Trim(secret, "\n")
	return secret, nil
}
