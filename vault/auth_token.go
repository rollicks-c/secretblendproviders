package vault

import (
	"github.com/hashicorp/vault/api"
)

type tokenProviderDirect struct {
	token string
}

func (td tokenProviderDirect) Authenticate(apiClient *api.Client) (*api.SecretAuth, error) {
	auth := &api.SecretAuth{
		ClientToken:   td.token,
		LeaseDuration: 60 * 60,
	}
	return auth, nil
}
