package cloudvolumes

import "github.com/chnsz/golangsdk"

const resourcePath = "cloudvolumes"
const resourceVolumePath = "volumes"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func retypeURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourceVolumePath, id, "retype")
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "action")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "detail")
}
