package clusters

import "github.com/chnsz/golangsdk"

const rootPath = "clusters"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL(rootPath, clusterId)
}
