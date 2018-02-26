package extensions

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Get retrieves information for a specific extension using its alias.
func Get(c *golangsdk.ServiceClient, alias string) (r GetResult) {
	_, r.Err = c.Get(ExtensionURL(c, alias), &r.Body, nil)
	return
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *golangsdk.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, ListExtensionURL(c), func(r pagination.PageResult) pagination.Page {
		return ExtensionPage{pagination.SinglePageBase(r)}
	})
}
