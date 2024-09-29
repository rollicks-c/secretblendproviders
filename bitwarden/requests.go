package bitwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type genericResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func (c Client) doRequest(method, ep string, body io.Reader) ([]byte, error) {

	// prep request
	url := fmt.Sprint(c.apiServer.getURL(), ep)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// do request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// done
	return payload, nil
}

func (c Client) doTypedRequest(method, ep string, body interface{}, response interface{}) error {

	// prep payload
	var payloadRaw io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err

		}
		payloadRaw = bytes.NewReader(payload)
	}

	// do request
	res, err := c.doRequest(method, ep, payloadRaw)
	if err != nil {
		return err
	}

	// parse response
	err = json.Unmarshal(res, response)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	return nil
}
