package bitwarden

import "time"

type totpResponse struct {
	Success bool     `json:"success"`
	Data    TOTPData `json:"data"`
}

type TOTPData struct {
	Object string `json:"object"`
	Data   string `json:"data"`
}

type itemResponse struct {
	Success bool     `json:"success"`
	Data    ItemData `json:"data"`
}

type listResponse struct {
	Success bool     `json:"success"`
	Data    listData `json:"data"`
}

type listData struct {
	Data    []ItemData `json:"data"`
}

type ItemData struct {
	PasswordHistory interface{}   `json:"passwordHistory"`
	RevisionDate    time.Time     `json:"revisionDate"`
	CreationDate    time.Time     `json:"creationDate"`
	DeletedDate     interface{}   `json:"deletedDate"`
	Object          string        `json:"object"`
	Id              string        `json:"id"`
	OrganizationId  interface{}   `json:"organizationId"`
	FolderId        interface{}   `json:"folderId"`
	Type            int           `json:"type"`
	Reprompt        int           `json:"reprompt"`
	Name            string        `json:"name"`
	Notes           interface{}   `json:"notes"`
	Favorite        bool          `json:"favorite"`
	Login           itemLogin     `json:"login"`
	CollectionIds   []interface{} `json:"collectionIds"`
}

type itemLogin struct {
	Fido2Credentials     []interface{} `json:"fido2Credentials"`
	Uris                 []itemURL     `json:"uris"`
	Username             string        `json:"username"`
	Password             string        `json:"password"`
	Totp                 string        `json:"totp"`
	PasswordRevisionDate interface{}   `json:"passwordRevisionDate"`
}

type itemURL struct {
	Match interface{} `json:"match"`
	Uri   string      `json:"uri"`
}
