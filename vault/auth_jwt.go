package vault

import (
	"github.com/hashicorp/vault/api"
)

type tokenProviderJWT struct {
	jwt      string
	authPath string
	role     string
	token    string
}

func (tp *tokenProviderJWT) Authenticate(apiClient *api.Client) (*api.SecretAuth, error) {
	args := map[string]any{
		"jwt":  tp.jwt,
		"role": tp.role,
	}
	secret, err := apiClient.Logical().Write(tp.authPath, args)
	if err != nil {
		return nil, err
	}
	return secret.Auth, nil
}
