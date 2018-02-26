package groups

import (
	"log"

	"github.com/huaweicloud/golangsdk"
)

const resourcePath = "scaling_group"

func createURL(c *golangsdk.ServiceClient) string {
	ur := c.ServiceURL(c.ProjectID, resourcePath)
	log.Printf("[DEBUG] Create URL is: %#v", ur)
	return ur
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}

func enableURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "action")
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}
