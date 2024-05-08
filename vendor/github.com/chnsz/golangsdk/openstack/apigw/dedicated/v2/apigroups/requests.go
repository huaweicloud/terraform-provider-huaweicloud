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

// AssociateDomainOpts is the structure that used to bind domain name to specified API group.
type AssociateDomainOpts struct {
	// The dedicated instance ID.
	InstanceId string `json:"-" requires:"true"`
	// The API group ID.
	GroupId string `json:"-" requires:"true"`
	// Custom domain name.
	// The valid length is limited from 0 to 255 characters and must comply with the domian name specifications.
	UrlDomain string `json:"url_domain" requires:"true"`
	// The minimum TLS version that can be used to access the domain name, the default value is TLSv1.2.
	// The valid value are as follows:
	// + TLSv1.1.
	// + TLSv1.2.
	MinSSLVersion string `json:"min_ssl_version,omitempty"`
	// Whether to enable redirection from HTTP to HTTPS, the default value is false.
	IsHttpRedirectToHttps bool `json:"is_http_redirect_to_https,omitempty"`
}

// AssociateDomain is a method that used to bind domain name to specified API group.
func AssociateDomain(client *golangsdk.ServiceClient, opts AssociateDomainOpts) (*AssociateDoaminResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r AssociateDoaminResp
	_, err = client.Post(associateDomainURL(client, opts.InstanceId, opts.GroupId), b, &r, nil)
	return &r, err
}

// AssociateDomain is a method that used to unbind domain name to specified API group.
func DisAssociateDomain(client *golangsdk.ServiceClient, intanceId string, groupId string, domainId string) error {
	_, err := client.Delete(disAssociateDomainURL(client, intanceId, groupId, domainId), nil)
	return err
}

// UpdateDomainAccessEnabledOpts is the structure that whether to use the dubugging domain name access the APIs.
type UpdateDomainAccessEnabledOpts struct {
	// The ID of the instance to which the group belongs.
	InstanceId string `json:"-" required:"true"`
	// The ID of the group.
	GroupId string `json:"-" required:"true"`
	// Whether to use the debugging domain name to access the APIs within the group.
	// Defalut value is true.
	SlDomainAccessEnabled *bool `json:"sl_domain_access_enabled" required:"true"`
}

// UpdateDomainAccessEnabled is a method used to control whether the APIs in the group can be accessed
// through the dubugging domain.
func UpdateDomainAccessEnabled(c *golangsdk.ServiceClient, opts UpdateDomainAccessEnabledOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(domainAccessEnabledURL(c, opts.InstanceId, opts.GroupId), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return err
}
