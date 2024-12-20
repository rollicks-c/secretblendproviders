package vault

import (
	"github.com/hashicorp/vault/api"
)

type tokenProviderAppRole struct {
	roleID   string
	secretID string
	authPath string
}

func (tar tokenProviderAppRole) Authenticate(apiClient *api.Client) (*api.SecretAuth, error) {

	// login and get a token
	tokenData, err := apiClient.Logical().Write(tar.authPath, map[string]interface{}{
		"role_id":   tar.roleID,
		"secret_id": tar.secretID,
	})
	if err != nil {
		return nil, err
	}

	// extract token and use
	return tokenData.Auth, nil
}
