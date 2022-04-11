package groups

import (
	"fmt"

	"github.com/chnsz/golangsdk"
)

const resourcePath = "scaling_group"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func forceDeleteURL(c *golangsdk.ServiceClient, id string) string {
	url := c.ServiceURL(resourcePath, id)
	return fmt.Sprintf("%s?force_delete=yes", url)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func enableURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "action")
}
