package tokens

import "github.com/chnsz/golangsdk"

func tokenURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
