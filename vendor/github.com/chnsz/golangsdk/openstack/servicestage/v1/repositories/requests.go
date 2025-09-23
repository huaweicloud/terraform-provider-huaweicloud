package repositories

import (
	"github.com/chnsz/golangsdk"
)

// PwdAuthOpts is the structure required by the CreatePwdAuth method to create the authorization.
type PwdAuthOpts struct {
	// Specified the authorization name.
	Name string `json:"name" required:"true"`
	// Specified the base64 encoded token, before encoding, the format is '{account name}:{password}'.
	Token string `json:"token" required:"true"`
}

// CreatePwdAuth is a method to create password authorization for a Git repository.
func CreatePwdAuth(c *golangsdk.ServiceClient, rType string, opts PwdAuthOpts) (*Authorization, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(passwordURL(c, rType), b, &rst.Body, nil)
	if err == nil {
		var r Authorization
		rst.ExtractIntoStructPtr(&r, "authorization")
		return &r, nil
	}
	return nil, err
}

// PersonalAuthOpts is the structure required by the CreatePersonalAuth method to create the authorization.
type PersonalAuthOpts struct {
	// Specified the authorization name.
	Name string `json:"name" required:"true"`
	// Specified the repository token.
	Token string `json:"token" required:"true"`
	// Specified the repository address.
	Host string `json:"host,omitempty"`
}

// CreatePersonalAuth is a method to create the personal access token authorization.
func CreatePersonalAuth(c *golangsdk.ServiceClient, rType string, opts PersonalAuthOpts) (*Authorization, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(personalURL(c, rType), b, &rst.Body, nil)
	if err == nil {
		var r Authorization
		rst.ExtractIntoStructPtr(&r, "authorization")
		return &r, nil
	}
	return nil, err
}

// List is a method to obtain the authorization list of the repositories.
func List(c *golangsdk.ServiceClient) ([]Authorization, error) {
	var rst golangsdk.Result
	_, err := c.Get(rootURL(c), &rst.Body, nil)
	if err == nil {
		var r []Authorization
		rst.ExtractIntoSlicePtr(&r, "authorizations")
		return r, nil
	}
	return nil, err
}

// Delete is a method to delete an existing authorization.
func Delete(c *golangsdk.ServiceClient, name string) error {
	_, err := c.Delete(resourceURL(c, name), nil)
	return err
}
