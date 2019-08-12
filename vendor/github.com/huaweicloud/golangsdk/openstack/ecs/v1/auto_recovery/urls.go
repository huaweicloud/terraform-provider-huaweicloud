package auto_recovery

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "cloudservers"
	resourcePath = "autorecovery"
)

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return updateURL(c, id)
}
