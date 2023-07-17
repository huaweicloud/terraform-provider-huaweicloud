package environments

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("envs")
}

func environmentURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("envs", id)
}
