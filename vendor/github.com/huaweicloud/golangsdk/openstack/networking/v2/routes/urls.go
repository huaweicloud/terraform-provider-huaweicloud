package routes

import "github.com/huaweicloud/golangsdk"

const resourcePath = "routes"
const rootPath = "vpc"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
