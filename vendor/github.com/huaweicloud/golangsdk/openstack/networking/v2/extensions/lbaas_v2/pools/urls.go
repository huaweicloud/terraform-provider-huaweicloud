package pools

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "lbaas"
	resourcePath = "pools"
	memberPath   = "members"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func memberRootURL(c *golangsdk.ServiceClient, poolId string) string {
	return c.ServiceURL(rootPath, resourcePath, poolId, memberPath)
}

func memberResourceURL(c *golangsdk.ServiceClient, poolID string, memeberID string) string {
	return c.ServiceURL(rootPath, resourcePath, poolID, memberPath, memeberID)
}
