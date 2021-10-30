package repositories

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOptsBuilder interface {
	ToRepositoryCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	// Name of the image repository
	Repository string `json:"repository" required:"true"`
	// Type of the image repository
	// The value can be app_server, linux, framework_app, database, lang, other, windows, arm
	Category string `json:"category,omitempty"`
	// Description of the image repository
	Description string `json:"description,omitempty"`
	// Whether the image repository is public
	IsPublic bool `json:"is_public"`
}

func (opts CreateOpts) ToRepositoryCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create new repository in the namespace(organization)
func Create(client *golangsdk.ServiceClient, namespace string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRepositoryCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, namespace), b, &r.Body, nil)
	return
}

type ListOptsBuilder interface {
	ToRepositoryListQuery() (string, error)
}

type ListOpts struct {
	// Namespace(Organization) name
	Namespace string `q:"namespace"`
	// Image repository name
	Name string `q:"name"`
	// Image repository category.
	// The value can be app_server, linux, framework_app, database, lang, other, windows, arm
	Category string `q:"category"`
	// Sorting by column.
	// You can set this parameter to name, updated_time, and tag_count.
	// The parameters OrderColumn and OrderType should always be used together.
	OrderColumn string `q:"order_column"`
	// Sorting type.
	// You can set this parameter to desc(descending sort) and asc(ascending sort).
	OrderType string `q:"order_type"`
	// offset 0 is a valid value
	Offset *int `q:"offset,omitempty"`
	Limit  int  `q:"limit,omitempty"`
}

const defaultLimit = 25

func (opts ListOpts) ToRepositoryListQuery() (string, error) {
	if opts.Limit == 0 && opts.Offset != nil {
		opts.Limit = defaultLimit
	}
	if opts.Limit != 0 && opts.Offset == nil {
		return "", fmt.Errorf("offset has to be defined if the limit is set")
	}
	if (opts.OrderColumn != "" && opts.OrderType == "") || (opts.OrderColumn == "" && opts.OrderType != "") {
		return "", fmt.Errorf("`OrderColumn` and `OrderType` should always be used together")
	}
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) (p pagination.Pager) {
	url := listURL(client)
	if opts != nil {
		q, err := opts.ToRepositoryListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += q
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RepositoryPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Get(client *golangsdk.ServiceClient, namespace, repository string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, namespace, repository), &r.Body, nil)
	return
}

type UpdateOptsBuilder interface {
	ToRepositoryUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"is_public"`
}

func (opts UpdateOpts) ToRepositoryUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Update(client *golangsdk.ServiceClient, namespace, repository string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRepositoryUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(resourceURL(client, namespace, repository), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, namespace, repository string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, namespace, repository), nil)
	return
}
