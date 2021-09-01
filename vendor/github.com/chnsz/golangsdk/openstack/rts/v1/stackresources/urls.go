package stackresources

import "github.com/chnsz/golangsdk"

func listURL(c *golangsdk.ServiceClient, stackName string) string {
	return c.ServiceURL("stacks", stackName, "resources")
}
