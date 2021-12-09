package mappings

import "github.com/chnsz/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "mappings", id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "mappings")
}
