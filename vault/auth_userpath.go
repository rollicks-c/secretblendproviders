package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

type tokenProviderUserPass struct {
	username string
	password string
	authPath string
}

func (up tokenProviderUserPass) Authenticate(apiClient *api.Client) (*api.SecretAuth, error) {

	// login and get a token
	authPath := fmt.Sprintf("%s/%s", up.authPath, up.username)
	tokenData, err := apiClient.Logical().Write(authPath, map[string]interface{}{
		"password": up.password,
	})
	if err != nil {
		return nil, err
	}

	// extract token and use
	return tokenData.Auth, nil
}
