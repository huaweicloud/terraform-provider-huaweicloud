package stacktemplates

import "github.com/chnsz/golangsdk"

func getURL(c *golangsdk.ServiceClient, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "template")
}
