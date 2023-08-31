package appauths

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure that used to authorize APIs to allow applications access.
type CreateOpts struct {
	// The Dedicated instance ID.
	InstanceId string `json:"-" required:"true"`
	// The ID of the environment in which the apps will be authorized.
	EnvId string `json:"env_id" required:"true"`
	// The ID list of the applications authorized to access the APIs.
	AppIds []string `json:"app_ids" required:"true"`
	// The authorized API IDs.
	ApiIds []string `json:"api_ids" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to authorize APIs to allow applications access using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) ([]Authorization, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = c.Post(rootURL(c, opts.InstanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Auths, err
}

// ListOpts is the structure used to querying authorized and unauthorized API information.
type ListOpts struct {
	// The instnace ID to which the application and APIs belong.
	InstanceId string `json:"-" required:"true"`
	// The application ID.
	AppId string `q:"app_id" required:"true"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// The authorized API ID.
	ApiId string `q:"app_id"`
	// The authorized API name.
	ApiName string `q:"api_name"`
	// The ID of the API group to which the authorized APIs belong.
	GroupId string `q:"group_id"`
	// The name of the API group to which the authorized APIs belong.
	GroupName string `q:"group_name"`
	// Parameter name (only 'name' is supported) for exact matching.
	EnvId string `q:"env_id"`
}

// List is a method to obtain the authorized API list under a specified application using given parameters.
func ListAuthorized(c *golangsdk.ServiceClient, opts ListOpts) ([]ApiAuthInfo, error) {
	url := listAuthorizedURL(c, opts.InstanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := AuthorizedPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractAuthorizedApis(pages)
}

// List is a method to obtain the unauthorized API list under a specified application using given parameters.
func ListUnauthorized(c *golangsdk.ServiceClient, opts ListOpts) ([]ApiOutlineInfo, error) {
	url := listUnathorizedURL(c, opts.InstanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := UnauthorizedPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractUnauthorizedApis(pages)
}

// Delete is a method used to unauthorize API from specified application using given parameters.
func Delete(c *golangsdk.ServiceClient, instanceId, authId string) error {
	_, err := c.Delete(resourceURL(c, instanceId, authId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
