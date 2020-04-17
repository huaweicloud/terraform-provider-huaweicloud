package flavors

import "github.com/huaweicloud/golangsdk"

func baseURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("flavors")
}
