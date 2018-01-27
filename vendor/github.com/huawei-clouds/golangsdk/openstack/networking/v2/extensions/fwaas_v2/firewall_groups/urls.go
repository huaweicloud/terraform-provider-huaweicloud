package firewall_groups

import "github.com/huawei-clouds/golangsdk"

const (
	rootPath     = "fwaas"
	resourcePath = "firewall_groups"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
