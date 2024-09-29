package bitwarden

import (
	"github.com/rollicks-c/term"
)

func (c Client) init() error {

	// check if already unlocked
	isLocked, err := c.IsLocked()
	if err != nil {
		return err
	}
	if !isLocked {
		return nil

	}

	// unlock with password
	pw, err := term.PromptSecret("enter bitwarden password")
	if err != nil {
		return err
	}
	if err := c.Unlock(pw); err != nil {
		return err
	}

	return nil
}
