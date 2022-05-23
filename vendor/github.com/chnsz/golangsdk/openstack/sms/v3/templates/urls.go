package templates

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("vm/templates")
}

func templateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("vm/templates", id)
}
