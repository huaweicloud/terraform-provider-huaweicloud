package tasks

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("tasks")
}

func taskURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("tasks", id)
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("tasks", id, "action")
}
