package vault

import (
	"fmt"
	"strings"
)

func (c Client) verifyToken() error {
	return c.authManager.refreshToken()
}

func (c Client) loadSecret(path string) (map[string]interface{}, error) {

	// login to vault
	vt, err := c.authManager.getClient()
	if err != nil {
		return nil, err
	}

	// apply path options
	path = c.fixDataPathForV2(path)

	// retrieve secret
	res, err := vt.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	// not found
	if res == nil {
		return nil, nil
	}

	// unpack
	secretDataRaw, ok := res.Data["data"]
	if !ok {
		return nil, fmt.Errorf("invalid secret: %s", path)
	}
	if secretDataRaw == nil {
		err := fmt.Errorf("found empty secret: %s", path)
		return nil, err
	}
	secretData := secretDataRaw.(map[string]interface{})

	// found
	return secretData, nil

}

func (c Client) fixDataPathForV2(secretPath string) string {
	secretPath = strings.TrimPrefix(secretPath, "/")
	parts := strings.Split(secretPath, "/")
	parts = append([]string{parts[0], "data"}, parts[1:]...)
	secretPath = strings.Join(parts, "/")
	return secretPath
}

func (c Client) fixMetaPathForV2(secretPath string) string {
	secretPath = strings.TrimPrefix(secretPath, "/")
	parts := strings.Split(secretPath, "/")
	parts = append([]string{parts[0], "metadata"}, parts[1:]...)
	secretPath = strings.Join(parts, "/")
	return secretPath
}
