package stacktemplates

import "github.com/huaweicloud/golangsdk"

func getURL(c *golangsdk.ServiceClient, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "template")
}
