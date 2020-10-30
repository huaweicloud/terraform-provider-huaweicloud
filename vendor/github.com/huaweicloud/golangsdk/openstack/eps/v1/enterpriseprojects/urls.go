package enterpriseprojects

import "github.com/huaweicloud/golangsdk"

const resourcePath = "enterprise-projects"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
