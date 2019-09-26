package certificates

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "lbaas"
	resourcePath = "certificates"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
