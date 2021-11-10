package domains

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOptsBuilder interface {
	ToAccessDomainCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	AccessDomain string `json:"access_domain" required:"true"`
	// Currently, only the `read` permission is supported.
	Permit string `json:"permit" required:"true"`
	// End date of image sharing (UTC). When the value is set to `forever`,
	// the image will be permanently available for the domain.
	// The validity period is calculated by day. The shared images expire at 00:00:00 on the day after the end date.
	Deadline    string `json:"deadline" required:"true"`
	Description string `json:"description,omitempty"`
}

func (opts CreateOpts) ToAccessDomainCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Create(client *golangsdk.ServiceClient, namespace, repository string, opts CreateOptsBuilder) (r CreateResult) {
	url := rootURL(client, namespace, repository)
	b, err := opts.ToAccessDomainCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(url, b, &r.Body, nil)
	return
}

func Get(client *golangsdk.ServiceClient, namespace, repository, domain string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, namespace, repository, domain), &r.Body, nil)
	return
}

type UpdateOptsBuilder interface {
	ToAccessDomainUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts used for update operations
// For argument details see CreateOpts
type UpdateOpts struct {
	Permit      string  `json:"permit" required:"true"`
	Deadline    string  `json:"deadline" required:"true"`
	Description *string `json:"description,omitempty"`
}

func (opts UpdateOpts) ToAccessDomainUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Update(client *golangsdk.ServiceClient, namespace, repository, domain string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAccessDomainUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Patch(resourceURL(client, namespace, repository, domain), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, namespace, repository, domain string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, namespace, repository, domain), nil)
	return
}

func List(client *golangsdk.ServiceClient, namespace, repository string) (p pagination.Pager) {
	return pagination.NewPager(client, rootURL(client, namespace, repository), func(r pagination.PageResult) pagination.Page {
		return AccessDomainPage{SinglePageBase: pagination.SinglePageBase(r)}
	})
}
