package sources

import "github.com/chnsz/golangsdk"

// listURL /sources
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("sources")
}
