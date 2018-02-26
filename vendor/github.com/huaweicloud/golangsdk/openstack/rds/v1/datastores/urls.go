package datastores

import "github.com/huaweicloud/golangsdk"

func listURL(c *golangsdk.ServiceClient, dataStoreName string) string {
	return c.ServiceURL("datastores", dataStoreName, "versions")
}
