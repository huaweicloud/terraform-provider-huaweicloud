package routers

import "github.com/huaweicloud/golangsdk"

const resourcePath = "routers"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func addInterfaceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "add_router_interface")
}

func removeInterfaceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "remove_router_interface")
}
