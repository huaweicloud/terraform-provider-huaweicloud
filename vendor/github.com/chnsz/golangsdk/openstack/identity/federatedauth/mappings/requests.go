package mappings

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type MappingOption struct {
	Rules []MappingRule `json:"rules" required:"true"`
}

type MappingRule struct {
	Local  []LocalRule  `json:"local" required:"true"`
	Remote []RemoteRule `json:"remote" required:"true"`
}

type LocalRule struct {
	User  *LocalRuleVal `json:"user,omitempty"`
	Group *LocalRuleVal `json:"group,omitempty"`
}

type LocalRuleVal struct {
	Name string `json:"name"`
}

type RemoteRule struct {
	Type     string   `json:"type" required:"true"`
	AnyOneOf []string `json:"any_one_of,omitempty"`
	NotAnyOf []string `json:"not_any_of,omitempty"`
}

type responseBody struct {
	Mapping IdentityMapping `json:"mapping"`
}

func Create(c *golangsdk.ServiceClient, id string, opts MappingOption) (*IdentityMapping, error) {
	b, err := golangsdk.BuildRequestBody(opts, "mapping")
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
		return &r.Mapping, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, id string) (*IdentityMapping, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, id), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r responseBody
		rst.ExtractInto(&r)
		return &r.Mapping, nil
	}
	return nil, err
}

func List(client *golangsdk.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return IdentityMappingPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Update(c *golangsdk.ServiceClient, id string, opts MappingOption) (*IdentityMapping, error) {
	b, err := golangsdk.BuildRequestBody(opts, "mapping")
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
		return &r.Mapping, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, id string) error {
	_, err := c.Delete(resourceURL(c, id), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return err
}
