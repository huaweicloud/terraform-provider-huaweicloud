package domains

import "github.com/chnsz/golangsdk"

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("active-domains")
}
