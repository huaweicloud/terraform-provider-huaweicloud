package credentials

import (
	"github.com/chnsz/golangsdk"
)

const parentElement = "credential"

type ListOptsBuilder interface {
	ToCredentialListQuery() (string, error)
}

type ListOpts struct {
	UserID string `json:"user_id,omitempty"`
}

func (opts ListOpts) ToCredentialListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) (l ListResult) {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToCredentialListQuery()
		if err != nil {
			l.Err = err
			return
		}
		url += query
	}

	_, l.Err = client.Get(url, &l.Body, nil)
	return
}

type CreateOptsBuilder interface {
	ToCredentialCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	UserID      string `json:"user_id" required:"true"`
	Description string `json:"description,omitempty"`
}

func (opts CreateOpts) ToCredentialCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, parentElement)
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCredentialCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), &b, &r.Body, nil)
	return
}

type UpdateOptsBuilder interface {
	ToCredentialUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	Status      string `json:"status,omitempty"`
	Description string `json:"description"`
}

func (opts UpdateOpts) ToCredentialUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, parentElement)
}

func Update(client *golangsdk.ServiceClient, credentialID string, opts UpdateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCredentialUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, credentialID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Get(client *golangsdk.ServiceClient, credentialID string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, credentialID), &r.Body, nil)
	return
}

func Delete(client *golangsdk.ServiceClient, credentialID string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, credentialID), nil)
	return
}
