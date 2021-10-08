package flavors

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func ListNodeTypes(c *golangsdk.ServiceClient) (*NodeTypes, error) {
	var rst NodeTypes
	_, err := c.Get(listNodeTypesURL(c), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		return &rst, nil
	}
	return nil, err
}
