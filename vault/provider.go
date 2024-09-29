package vault

type Provider struct {
	client *Client
}

func (p Provider) LoadSecret(uri string) (string, error) {
	return p.client.LoadSecretValue(uri)
}
