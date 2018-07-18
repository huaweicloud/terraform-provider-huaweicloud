package stacktemplates

import "github.com/huaweicloud/golangsdk"

// Get retreives data for the given stack template.
func Get(c *golangsdk.ServiceClient, stackName, stackID string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, stackName, stackID), &r.Body, nil)
	return
}
