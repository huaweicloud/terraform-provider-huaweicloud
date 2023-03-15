package dnats

import "github.com/chnsz/golangsdk"

const rootPath = "private-nat/dnat-rules"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, ruleId string) string {
	return c.ServiceURL(rootPath, ruleId)
}
