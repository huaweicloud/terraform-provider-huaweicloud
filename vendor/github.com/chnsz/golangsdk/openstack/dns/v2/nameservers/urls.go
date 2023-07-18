package nameservers

import "github.com/chnsz/golangsdk"

func baseURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("nameservers")
}
