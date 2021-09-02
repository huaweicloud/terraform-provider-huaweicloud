package authorizers

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CustomAuthOpts is a struct which will be used to create a new custom authorizer or
// update an existing custom authorizer.
type CustomAuthOpts struct {
	// Custom authorizer name, which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, and underscores (_) are allowed.
	Name string `json:"name" required:"true"`
	// Custom authorizer type, which support 'FRONTEND' and 'BACKEND'.
	Type string `json:"type" required:"true"`
	// Authorizer type, and the value is 'FUNC'.
	AuthorizerType string `json:"authorizer_type" required:"true"`
	// Function URN.
	AuthorizerURI string `json:"authorizer_uri" required:"true"`
	// Indicates whether to send the body.
	IsBodySend *bool `json:"need_body,omitempty"`
	// Identity source.
	Identities []AuthCreateIdentitiesReq `json:"identities,omitempty"`
	// Maximum cache age. The maximum value is 3,600.
	// The maximum length of time that authentication results can be cached for.
	// A value of 0 means that results are not cached.
	TTL *int `json:"ttl,omitempty"`
	// User data.
	UserData *string `json:"user_data,omitempty"`
}

// AuthCreateIdentitiesReq is an object which will be build up a indentities list.
type AuthCreateIdentitiesReq struct {
	// Parameter name.
	Name string `json:"name" required:"true"`
	// Parameter location, which support 'HEADER' and 'QUERY'.
	Location string `json:"location" required:"true"`
	// Parameter verification expression. The default value is null, indicating that no verification is performed.
	Validation *string `json:"validation,omitempty"`
}

// CustomAuthOptsBuilder is an interface which to support request body build of
// the custom authorizer creation and updation.
type CustomAuthOptsBuilder interface {
	ToCustomAuthOptsMap() (map[string]interface{}, error)
}

// ToCustomAuthOptsMap is a method which to build a request body by the CustomAuthOpts.
func (opts CustomAuthOpts) ToCustomAuthOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method to create a new custom authorizer.
func Create(client *golangsdk.ServiceClient, instanceId string, opts CustomAuthOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToCustomAuthOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// Update is a method to update an existing custom authorizer.
func Update(client *golangsdk.ServiceClient, instanceId, authId string, opts CustomAuthOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToCustomAuthOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, authId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get is a method to obtain an existing custom authorizer.
func Get(client *golangsdk.ServiceClient, instanceId, authId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, authId), &r.Body, nil)
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// ID.
	ID string `q:"id"`
	// Name.
	Name string `q:"name"`
	// Custom authorizer type.
	Type string `q:"type"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
}

// ListOptsBuilder is an interface which to support request query build of
// the custom authorizer search.
type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

// ToListQuery is a method which to build a request query by the ListOpts.
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more custom authorizers for the instance according to the
// query parameters.
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
		return CustomAuthPage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete an existing custom authorizer.
func Delete(client *golangsdk.ServiceClient, instanceId, authId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, authId), nil)
	return
}
