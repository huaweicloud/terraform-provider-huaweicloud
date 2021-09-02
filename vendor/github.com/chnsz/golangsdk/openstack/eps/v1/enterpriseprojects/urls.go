package enterpriseprojects

import "github.com/chnsz/golangsdk"

const resourcePath = "enterprise-projects"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
