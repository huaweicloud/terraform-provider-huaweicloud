package endpoints

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// BatchOpts is the structure that used to add or remove permissions for endpoint service.
type BatchOpts struct {
	// The instance ID to which the endpoint service belong.
	InstanceId string `json:"-" required:"true"`
	// The permissions of endpoint service to add or remove.
	Permissions []string `json:"permissions" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// AddPermissions is a method used to batch add permissions to endpoint service using given parameters.
func AddPermissions(c *golangsdk.ServiceClient, opts BatchOpts) ([]string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r createResp
	_, err = c.Post(addURL(c, opts.InstanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Permissions, err
}

// ListOpts is the structure that used to query the permissions of endpoint service.
type ListOpts struct {
	// The instance ID to which the endpoint service belong.
	InstanceId string `json:"-" required:"true"`
	// The permission of endpoint service.
	Permission string `json:"permission"`
	// Offset value. The value must be a positive integer.
	Offset int `q:"offset"`
	// Number of records displayed per page.
	// The value must be a positive integer.
	Limit int `q:"limit"`
}

// ListPermissions is a method used to query the permissions of endpoint service with given parameters.
func ListPermissions(c *golangsdk.ServiceClient, opts ListOpts) ([]EndpointPermission, error) {
	url := listURL(c, opts.InstanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := PermissionPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractPermissions(pages)
}

// DeletePermissions is a method used to batch delete permissions of endpoint service using given parameters.
func DeletePermissions(c *golangsdk.ServiceClient, opts BatchOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(deleteURL(c, opts.InstanceId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
