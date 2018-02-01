package backendecs

import "github.com/huawei-clouds/golangsdk"

const (
	rootPath     = "elbaas"
	resourcePath = "listeners"
	memberPath   = "members"
)

func rootURL(c *golangsdk.ServiceClientExtension, lId string) string {
	return c.ServiceURL(c.ProjectID, rootPath, resourcePath, lId, memberPath)
}

func actionURL(c *golangsdk.ServiceClientExtension, lId string) string {
	return c.ServiceURL(c.ProjectID, rootPath, resourcePath, lId, memberPath, "action")
}
