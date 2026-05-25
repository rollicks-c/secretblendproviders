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
		// Login.Username takes precedence; if empty (e.g. Identity-type
		// item) fall back to Identity.Username so callers can address
		// the identity field by its canonical name.
		if item.Login.Username != "" {
			return item.Login.Username, nil
		}
		return item.Identity.Username, nil
	case "password":
		return item.Login.Password, nil
	// Card-type fields (item.Type == 3). Names match the JSON returned
	// by `bw get item` so callers can request `number`, `code`, etc.
	case "cardholderName":
		return item.Card.CardholderName, nil
	case "brand":
		return item.Card.Brand, nil
	case "number":
		return item.Card.Number, nil
	case "expMonth":
		return item.Card.ExpMonth, nil
	case "expYear":
		return item.Card.ExpYear, nil
	case "code":
		return item.Card.Code, nil
	// Identity-type fields (item.Type == 4).
	case "title":
		return item.Identity.Title, nil
	case "firstName":
		return item.Identity.FirstName, nil
	case "middleName":
		return item.Identity.MiddleName, nil
	case "lastName":
		return item.Identity.LastName, nil
	case "address1":
		return item.Identity.Address1, nil
	case "address2":
		return item.Identity.Address2, nil
	case "address3":
		return item.Identity.Address3, nil
	case "city":
		return item.Identity.City, nil
	case "state":
		return item.Identity.State, nil
	case "postalCode":
		return item.Identity.PostalCode, nil
	case "country":
		return item.Identity.Country, nil
	case "company":
		return item.Identity.Company, nil
	case "email":
		return item.Identity.Email, nil
	case "phone":
		return item.Identity.Phone, nil
	case "ssn":
		return item.Identity.SSN, nil
	case "passportNumber":
		return item.Identity.PassportNumber, nil
	case "licenseNumber":
		return item.Identity.LicenseNumber, nil
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
