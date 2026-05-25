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
	Message string   `json:"message"`
}

type listResponse struct {
	Success bool     `json:"success"`
	Data    listData `json:"data"`
}

type listData struct {
	Data []ItemData `json:"data"`
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
	Card            ItemCard      `json:"card"`
	Identity        ItemIdentity  `json:"identity"`
	Fields          []ItemField   `json:"fields"`
	CollectionIds   []interface{} `json:"collectionIds"`
}

// ItemCard mirrors the `card` object returned by `bw get item` for
// Card-type vault entries (Type=3). Field names are kept as Bitwarden
// reports them so callers can ask for `number`, `code`, `expMonth`, etc.
type ItemCard struct {
	CardholderName string `json:"cardholderName"`
	Brand          string `json:"brand"`
	Number         string `json:"number"`
	ExpMonth       string `json:"expMonth"`
	ExpYear        string `json:"expYear"`
	Code           string `json:"code"`
}

// ItemIdentity mirrors the `identity` object returned by `bw get item`
// for Identity-type vault entries (Type=4).
type ItemIdentity struct {
	Title          string `json:"title"`
	FirstName      string `json:"firstName"`
	MiddleName     string `json:"middleName"`
	LastName       string `json:"lastName"`
	Address1       string `json:"address1"`
	Address2       string `json:"address2"`
	Address3       string `json:"address3"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postalCode"`
	Country        string `json:"country"`
	Company        string `json:"company"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	SSN            string `json:"ssn"`
	Username       string `json:"username"`
	PassportNumber string `json:"passportNumber"`
	LicenseNumber  string `json:"licenseNumber"`
}

// ItemField is a Bitwarden custom field on an item.
// Type: 0=Text, 1=Hidden, 2=Boolean, 3=Linked.
type ItemField struct {
	Name     string      `json:"name"`
	Value    string      `json:"value"`
	Type     int         `json:"type"`
	LinkedId interface{} `json:"linkedId"`
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
