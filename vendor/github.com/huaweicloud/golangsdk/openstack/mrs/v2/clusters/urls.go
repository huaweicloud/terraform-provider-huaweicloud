package clusters

import "github.com/huaweicloud/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("clusters")
}
