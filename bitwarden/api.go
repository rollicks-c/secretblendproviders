package bitwarden

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rollicks-c/configcove"
	"github.com/rollicks-c/secretblend"
)

type Client struct {
	apiServer *apiServer
}

type Option func(c *Client)

func Register() error {

	// use defaults
	appID := "bitwarden"
	dataDir := configcove.ConfigDir(appID)

	// register client
	bw, err := NewClient(dataDir)
	if err != nil {
		return err
	}
	secretblend.AddProvider(bw, "bitwarden://")

	return nil
}

func NewClient(dataDir string, options ...Option) (*Client, error) {

	// start api server
	server := &apiServer{
		dataDir: dataDir,
	}
	if err := server.start(); err != nil {
		return nil, err
	}

	// create api client
	client := &Client{
		apiServer: server,
	}
	for _, opt := range options {
		opt(client)
	}

	// get ready to accept requests
	if err := client.init(); err != nil {
		return nil, err
	}

	return client, nil

}

func (c Client) LoadSecret(uri string) (string, error) {
	parts := strings.Split(uri, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid uri: %s", uri)
	}

	itemID := parts[0]
	item, err := c.GetItem(itemID)
	if err != nil {
		return "", err
	}

	keyExp := parts[1]
	switch keyExp {
	case "username":
		return item.Login.Username, nil
	case "password":
		return item.Login.Password, nil
	default:
		return "", fmt.Errorf("invalid key: %s", keyExp)
	}

}

func (c Client) Check() error {

	ep := "/object/fingerprint/me"
	res := genericResponse{}
	if err := c.doTypedRequest(http.MethodGet, ep, nil, &res); err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.Message)
	}
	return nil

}

func (c Client) IsLocked() (bool, error) {

	err := c.Check()
	if err == nil {
		return false, nil
	}
	if err.Error() == "Vault is locked." {
		return true, nil
	}
	return true, err

}

func (c Client) Unlock(password string) error {
	ep := "/unlock"
	type request struct {
		Password string `json:"password"`
	}
	req := request{
		Password: password,
	}
	res := genericResponse{}
	if err := c.doTypedRequest(http.MethodPost, ep, req, &res); err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.Message)
	}
	return nil
}

func (c Client) Find(exp string) ([]ItemData, error) {

	ep := fmt.Sprintf("/list/object/items?search=%s", exp)
	res := listResponse{}

	if err := c.doTypedRequest(http.MethodGet, ep, nil, &res); err != nil {
		return nil, err
	}
	if !res.Success {
		return nil, fmt.Errorf("failed to get item: %v", res)
	}
	return res.Data.Data, nil

}

func (c Client) GetItem(id string) (ItemData, error) {

	ep := fmt.Sprintf("/object/item/%s", id)
	res := itemResponse{}
	if err := c.doTypedRequest(http.MethodGet, ep, nil, &res); err != nil {
		return ItemData{}, err
	}
	if !res.Success {
		return ItemData{}, fmt.Errorf("failed to get item: %v", res.Message)
	}
	return res.Data, nil

}

func (c Client) GetTOTP(id string) (TOTPData, error) {

	ep := fmt.Sprintf("/object/totp/%s", id)
	res := totpResponse{}
	if err := c.doTypedRequest(http.MethodGet, ep, nil, &res); err != nil {
		return TOTPData{}, err
	}
	if !res.Success {
		return TOTPData{}, fmt.Errorf("failed to get totp: %v", res)
	}
	return res.Data, nil

}

func (c Client) Sync() error {

	ep := "/sync"
	res := genericResponse{}

	if err := c.doTypedRequest(http.MethodPost, ep, nil, &res); err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("failed to sync bw vault: %v", res)
	}
	return nil

}
