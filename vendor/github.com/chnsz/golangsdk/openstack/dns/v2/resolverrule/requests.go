package resolverrule

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	Name        string      `json:"name" required:"true"`
	DomainName  string      `json:"domain_name"  required:"true"`
	EndpointID  string      `json:"endpoint_id"  required:"true"`
	IPAddresses []IPAddress `json:"ipaddresses"  required:"true"`
}

type IPAddress struct {
	IP string `json:"ip,omitempty"`
}

type UpdateOpts struct {
	Name        string      `json:"name,omitempty"`
	IPAddresses []IPAddress `json:"ipaddresses,omitempty"`
}

type ListOpts struct {
	Limit      int    `q:"limit"`
	Offset     int    `q:"offset"`
	DomainName string `q:"domain_name"`
}

func Create(client *golangsdk.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(baseURL(client), b, &r.Body, nil)
	return
}

func Update(client *golangsdk.ServiceClient, resolverRuleID string, opts UpdateOpts) (r UpdateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(resourceURL(client, resolverRuleID), b, &r.Body, nil)
	return
}

func Get(client *golangsdk.ServiceClient, resolverRuleID string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, resolverRuleID), &r.Body, nil)
	return
}

func List(client *golangsdk.ServiceClient, opts *ListOpts) (r ListResult) {
	url := baseURL(client)
	if opts != nil {
		query, err := golangsdk.BuildQueryString(opts)
		if err != nil {
			r.Err = err
			return
		}
		url += query.String()
	}

	r.Body, r.Err = pagination.ListAllItems(client, pagination.Offset, url, nil)
	return
}

func Delete(client *golangsdk.ServiceClient, resolverRuleID string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, resolverRuleID), nil)
	return
}
