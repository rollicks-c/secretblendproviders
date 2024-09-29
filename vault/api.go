package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"path"
)

type Client struct {
	authManager *authManager
}

type Option func(c *Client)

func NewClient(addr string, options ...Option) (*Client, error) {

	// create api client
	vtClient, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}

	client := &Client{
		authManager: newAuthManager(vtClient),
	}
	for _, opt := range options {
		opt(client)
	}

	if client.authManager.tokenProvider == nil {
		return nil, fmt.Errorf("no token provider is set")
	}

	if err := client.verifyToken(); err != nil {
		return nil, err
	}

	return client, nil

}

func WithAppRole(roleID string, secretID string) Option {
	tp := tokenProviderAppRole{
		roleID:   roleID,
		secretID: secretID,
		authPath: "auth/approle/login",
	}
	return func(c *Client) {
		c.authManager.tokenProvider = tp
	}
}

func WithJWT(authPath, role, jwt string) Option {
	tp := &tokenProviderJWT{
		jwt:      jwt,
		role:     role,
		authPath: authPath,
	}
	return func(c *Client) {
		c.authManager.tokenProvider = tp
	}
}

func WithToken(token string) Option {
	tp := tokenProviderDirect{
		token: token,
	}
	return func(c *Client) {
		if token == "" {
			return
		}
		c.authManager.tokenProvider = tp
	}
}

func (c Client) LoadSecret(uri string) (string, error) {

	vtPath := path.Dir(uri)
	key := path.Base(uri)

	data, err := c.loadSecret(vtPath)
	if err != nil {
		return "", err
	}

	secret, ok := data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return secret.(string), nil
}

func (c Client) ListSecret(path string) ([]string, error) {

	// login to vault
	vt, err := c.authManager.getClient()
	if err != nil {
		return nil, err
	}

	// apply path options
	path = c.fixMetaPathForV2(path)

	// retrieve secret
	res, err := vt.Logical().List(path)
	if err != nil {
		return nil, err
	}

	// not found
	if res == nil {
		return nil, fmt.Errorf("no value found at [%s]", path)
	}

	// unpack
	DataRaw, ok := res.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("invalid secret: %s", path)
	}
	data := DataRaw.([]interface{})
	keys := make([]string, 0, len(data))
	for _, k := range data {
		keys = append(keys, k.(string))
	}

	// found
	return keys, nil

}

func (c Client) WriteSecret(path string, data map[string]interface{}) error {

	// login to vault
	vt, err := c.authManager.getClient()
	if err != nil {
		return err
	}

	// pack secret
	payload := map[string]interface{}{
		"data": data,
	}

	// apply path options
	path = c.fixDataPathForV2(path)

	// write secret
	if _, err := vt.Logical().Write(path, payload); err != nil {
		return err
	}

	return nil

}

func (c Client) DeleteSecret(path string) error {

	// login to vault
	vt, err := c.authManager.getClient()
	if err != nil {
		return err
	}

	// apply path options
	//path = c.fixDataPathForV2(path)
	path = c.fixMetaPathForV2(path)

	// remove secret
	if _, err := vt.Logical().Delete(path); err != nil {
		return err
	}

	return nil

}

func (c Client) ReadValue(path, field string) (interface{}, error) {

	vt, err := c.authManager.getClient()
	if err != nil {
		return nil, err
	}

	// retrieve secret
	res, err := vt.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	// not found
	if res == nil {
		return nil, fmt.Errorf("no value found at [%s]", path)
	}

	// extract
	value, ok := res.Data[field]
	if !ok {
		return nil, fmt.Errorf("field %s not found at %s", field, path)
	}

	// found
	return value, nil

}
