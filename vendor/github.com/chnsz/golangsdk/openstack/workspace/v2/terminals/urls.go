package terminals

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("terminals/binding-desktops")
}

func deleteURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("terminals/binding-desktops/batch-delete")
}

func configURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("terminals/binding-desktops/config")
}
