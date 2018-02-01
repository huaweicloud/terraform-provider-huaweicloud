package healthcheck

import "github.com/huawei-clouds/golangsdk"

const (
	rootPath     = "elbaas"
	resourcePath = "healthcheck"
)

func rootURL(c *golangsdk.ServiceClientExtension) string {
	return c.ServiceURL(c.ProjectID, rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClientExtension, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, resourcePath, id)
}
