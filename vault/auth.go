package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"sync"
	"time"
)

type tokenProvider interface {
	Authenticate(apiClient *api.Client) (*api.SecretAuth, error)
}

type authManager struct {
	tokenProvider tokenProvider
	token         string
	tokenExp      *time.Time
	apiClient     *api.Client
	lock          sync.Mutex
}

func newAuthManager(apiClient *api.Client) *authManager {
	return &authManager{
		apiClient: apiClient,
		lock:      sync.Mutex{},
	}
}

func (am *authManager) getClient() (*api.Client, error) {

	am.lock.Lock()
	defer am.lock.Unlock()

	if err := am.refreshToken(); err != nil {
		return nil, err
	}
	return am.apiClient, nil
}

func (am *authManager) refreshToken() error {

	if err := am.validateToken(); err == nil {
		return nil
	}

	if err := am.renewToken(); err != nil {
		return err
	}

	return nil
}

func (am *authManager) validateToken() error {

	if am.token == "" {
		return fmt.Errorf("no token set")
	}
	if am.tokenExp == nil {
		return fmt.Errorf("no token meta dataa set")
	}
	if am.tokenExp.Before(time.Now()) {
		return fmt.Errorf("token has expired")
	}

	return nil
}

func (am *authManager) renewToken() error {

	// reset
	am.token = ""
	am.tokenExp = nil

	// gather token
	auth, err := am.tokenProvider.Authenticate(am.apiClient)
	if err != nil {
		return err
	}
	am.token = auth.ClientToken
	am.apiClient.SetToken(auth.ClientToken)

	// gather meta data
	tokenMeta, err := am.apiClient.Auth().Token().LookupSelf()
	if err != nil {
		return err
	}
	ttl, err := tokenMeta.TokenTTL()
	if err != nil {
		return err
	}
	expire := time.Now().Add(ttl).Add(time.Second * -5)

	expire = time.Now().Add(time.Duration(auth.LeaseDuration) * time.Second).Add(time.Second * -5)

	// use token
	am.tokenExp = &expire
	return nil
}

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
