package flavors

import "github.com/chnsz/golangsdk"

func baseURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("flavors")
}
