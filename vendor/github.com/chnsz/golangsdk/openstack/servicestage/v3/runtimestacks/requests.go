package runtimestacks

import (
	"github.com/chnsz/golangsdk"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// List is a method to obtain the details of a specified .
func List(c *golangsdk.ServiceClient) ([]RuntimeStack, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(nil)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}
	var r []RuntimeStack
	rst.ExtractIntoSlicePtr(&r, "runtimestacks")
	return r, nil
}
