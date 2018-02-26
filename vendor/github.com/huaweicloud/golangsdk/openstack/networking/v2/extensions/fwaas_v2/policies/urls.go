package policies

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "fwaas"
	resourcePath = "firewall_policies"
	insertPath   = "insert_rule"
	removePath   = "remove_rule"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func insertURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, insertPath)
}

func removeURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, removePath)
}
