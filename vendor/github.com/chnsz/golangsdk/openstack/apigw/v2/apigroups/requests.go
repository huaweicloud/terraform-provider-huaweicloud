package apigroups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// GroupOpts allows to create a group or update a existing group using given parameters.
type GroupOpts struct {
	// API group name, which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// Description of the API group, which can contain a maximum of 255 characters,
	// and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description *string `json:"remark,omitempty"`
}

type CreateOptsBuilder interface {
	ToCreateOptsMap() (map[string]interface{}, error)
}

func (opts GroupOpts) ToCreateOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create function that create a new group.
func Create(client *golangsdk.ServiceClient, instanceId string, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToCreateOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// Update is a method by which to create function that udpate a existing group.
func Update(client *golangsdk.ServiceClient, instanceId, groupId string, opts CreateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToCreateOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, groupId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get is a method to obtain the specified group according to the instanceId and appId.
func Get(client *golangsdk.ServiceClient, instanceId, groupId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, groupId), &r.Body, nil)
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// API group ID.
	Id string `q:"id"`
	// API group name.
	Name string `q:"name"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// Parameter name for exact matching. Only API group names are supported.
	PreciseSearch string `q:"precise_search"`
}

type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more groups according to the query parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, instanceId)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GroupPage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete an existing group.
func Delete(client *golangsdk.ServiceClient, instanceId, groupId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, groupId), nil)
	return
}
