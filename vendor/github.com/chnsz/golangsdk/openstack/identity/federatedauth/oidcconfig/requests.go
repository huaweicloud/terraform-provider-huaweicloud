package oidcconfig

import "github.com/chnsz/golangsdk"

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateOpts struct {
	AccessMode            string `json:"access_mode" required:"true"`
	IdpURL                string `json:"idp_url" required:"true"`
	ClientID              string `json:"client_id" required:"true"`
	SigningKey            string `json:"signing_key" required:"true"`
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	Scope                 string `json:"scope,omitempty"`
	ResponseType          string `json:"response_type,omitempty"`
	ResponseMode          string `json:"response_mode,omitempty"`
}

type responseBody struct {
	Config OpenIDConnectConfig `json:"openid_connect_config"`
}

func Create(c *golangsdk.ServiceClient, idpID string, opts CreateOpts) (*OpenIDConnectConfig, error) {
	b, err := golangsdk.BuildRequestBody(opts, "openid_connect_config")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(resourceURL(c, idpID), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r responseBody
		rst.ExtractInto(&r)
		return &r.Config, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, idpID string) (*OpenIDConnectConfig, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, idpID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	if err == nil {
		var r responseBody
		rst.ExtractInto(&r)
		return &r.Config, nil
	}
	return nil, err
}

type UpdateOpenIDConnectConfigOpts struct {
	AccessMode            string `json:"access_mode" required:"true"`
	IdpURL                string `json:"idp_url" required:"true"`
	ClientID              string `json:"client_id" required:"true"`
	SigningKey            string `json:"signing_key" required:"true"`
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	Scope                 string `json:"scope,omitempty"`
	ResponseType          string `json:"response_type,omitempty"`
	ResponseMode          string `json:"response_mode,omitempty"`
}

func Update(c *golangsdk.ServiceClient, idpID string, opts UpdateOpenIDConnectConfigOpts) (*OpenIDConnectConfig, error) {
	b, err := golangsdk.BuildRequestBody(opts, "openid_connect_config")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, idpID), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r responseBody
		rst.ExtractInto(&r)
		return &r.Config, nil
	}
	return nil, err
}
