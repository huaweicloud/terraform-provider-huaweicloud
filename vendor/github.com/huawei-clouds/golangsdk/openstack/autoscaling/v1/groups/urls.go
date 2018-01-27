package groups

import (
	"github.com/huawei-clouds/golangsdk"
	"log"
)

const resourcePath = "scaling_group"

func createURL(c *golangsdk.ServiceClientExtension) string {
	ur := c.ServiceURL(c.ProjectID, resourcePath)
	log.Printf("[DEBUG] Create URL is: %#v", ur)
	return ur
}

func deleteURL(c *golangsdk.ServiceClientExtension, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func getURL(c *golangsdk.ServiceClientExtension, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func listURL(c *golangsdk.ServiceClientExtension) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}

func enableURL(c *golangsdk.ServiceClientExtension, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "action")
}

func updateURL(c *golangsdk.ServiceClientExtension, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}
