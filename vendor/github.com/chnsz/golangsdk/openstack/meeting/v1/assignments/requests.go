package assignments

import (
	"fmt"

	"github.com/chnsz/golangsdk"
)

// CreateOpts is the structure required by the Create method to assign the administrator role for an account.
type CreateOpts struct {
	// Account type.
	//   0: HUAWEI CLOUD conference account. Used for account/password authentication.
	//   1: Third-party User ID, used for App ID authentication.
	// default 0
	AccountType int `json:"-"`
	// account account.
	// If it is an account/password authentication method, it refers to the HUAWEI CLOUD conference account.
	// If it is the App ID authentication method, it refers to the third-party User ID.
	Account string `json:"account" required:"true"`
	// Authorization token.
	Token string `json:"-" required:"true"`
}

// QueryOpts is the structure used to build the query path.
type QueryOpts struct {
	// Account type.
	//   0: HUAWEI CLOUD conference account. Used for account/password authentication.
	//   1: Third-party User ID, used for App ID authentication.
	// default 0
	AccountType int `q:"accountType"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to assign the administrator role for an account.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) error {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(QueryOpts{
		AccountType: opts.AccountType,
	})
	if err != nil {
		return err
	}
	url += query.String()

	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}
	_, err = c.Post(url, b, nil, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return err
}

// GetOpts is the structure used to build the query path and authorization.
type GetOpts struct {
	// Account type.
	//   0: HUAWEI CLOUD conference account. Used for account/password authentication.
	//   1: Third-party User ID, used for App ID authentication.
	// default 0
	AccountType int `q:"accountType"`
	// account account.
	// If it is an account/password authentication method, it refers to the HUAWEI CLOUD conference account.
	// If it is the App ID authentication method, it refers to the third-party User ID.
	Account string `json:"-" required:"true"`
	// Authorization token.
	Token string `json:"-" required:"true"`
}

// Get is a method to query the details of the role assignment.
func Get(c *golangsdk.ServiceClient, opts GetOpts) (*Administrator, error) {
	url := resourceURL(c, opts.Account)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r Administrator
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return &r, err
}

// DeleteOpts is the structure used to build the query path and authorization.
type DeleteOpts struct {
	// Account type.
	//   0: HUAWEI CLOUD conference account. Used for account/password authentication.
	//   1: Third-party User ID, used for App ID authentication.
	// default 0
	AccountType int `q:"accountType"`
	// Authorization token.
	Token string `json:"-"`
	// List of user account.
	Accounts []string `json:"-"`
}

// BatchDelete is a method to revoke administrator roles for all accounts.
func BatchDelete(c *golangsdk.ServiceClient, opts DeleteOpts) error {
	url := deleteURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return err
	}
	url += query.String()

	if opts.Token == "" {
		return fmt.Errorf("The authorization token must be supported.")
	}

	_, err = c.Post(url, opts.Accounts, nil, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return err
}
