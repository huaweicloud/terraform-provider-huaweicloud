package protectiongroups

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("server-groups")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("server-groups", id)
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("server-groups", id, "action")
}
