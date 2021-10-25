package providers

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateProviderOpts struct {
	SsoType     string `json:"sso_type,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
}

type responseBody struct {
	IdentityProvider Provider `json:"identity_provider"`
}

func Create(c *golangsdk.ServiceClient, id string, opts CreateProviderOpts) (*Provider, error) {
	b, err := golangsdk.BuildRequestBody(opts, "identity_provider")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, id), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r responseBody
		rst.ExtractInto(&r)
		return &r.IdentityProvider, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, id string) (*Provider, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, id), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r responseBody
		rst.ExtractInto(&r)
		return &r.IdentityProvider, nil
	}
	return nil, err
}

type UpdateOpts struct {
	Description *string `json:"description,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (*Provider, error) {
	b, err := golangsdk.BuildRequestBody(opts, "identity_provider")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Patch(resourceURL(c, id), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r responseBody
		rst.ExtractInto(&r)
		return &r.IdentityProvider, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, regionID string) error {
	_, err := c.Delete(resourceURL(c, regionID), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return err
}
