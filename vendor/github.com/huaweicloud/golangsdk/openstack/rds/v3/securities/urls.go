package securities

import "github.com/huaweicloud/golangsdk"

func rootURL(c *golangsdk.ServiceClient, instanceId, path string) string {
	return c.ServiceURL("instances", instanceId, path)
}
