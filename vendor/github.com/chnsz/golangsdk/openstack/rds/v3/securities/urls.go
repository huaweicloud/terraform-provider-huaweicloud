package securities

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, instanceId, path string) string {
	return c.ServiceURL("instances", instanceId, path)
}
