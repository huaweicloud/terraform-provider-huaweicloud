package enterpriseprojects

import "github.com/chnsz/golangsdk"

const resourcePath = "enterprise-projects"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "action")
}
