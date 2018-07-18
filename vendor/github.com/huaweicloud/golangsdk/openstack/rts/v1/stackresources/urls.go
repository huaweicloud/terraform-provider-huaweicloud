package stackresources

import "github.com/huaweicloud/golangsdk"

func listURL(c *golangsdk.ServiceClient, stackName string) string {
	return c.ServiceURL("stacks", stackName, "resources")
}
