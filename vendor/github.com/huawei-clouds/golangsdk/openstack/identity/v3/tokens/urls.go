package tokens

import "github.com/huawei-clouds/golangsdk"

func tokenURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
