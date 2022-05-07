package metadata

import (
	"github.com/chnsz/golangsdk"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// ListRuntimes is a method to obtain all component runtimes.
func ListRuntimes(c *golangsdk.ServiceClient) ([]Runtime, error) {
	var rst golangsdk.Result
	_, err := c.Get(runtimeURL(c), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	var r []Runtime
	err = rst.ExtractIntoSlicePtr(&r, "runtimes")
	return r, err
}

// ListFlavors is a method to obtain all application flavors.
func ListFlavors(c *golangsdk.ServiceClient) ([]Flavor, error) {
	var rst golangsdk.Result
	_, err := c.Get(flavorURL(c), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	var r []Flavor
	err = rst.ExtractIntoSlicePtr(&r, "flavors")
	return r, err
}
