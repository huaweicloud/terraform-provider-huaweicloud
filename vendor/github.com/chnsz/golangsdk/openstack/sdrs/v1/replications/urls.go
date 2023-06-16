package replications

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("replications")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("replications", id)
}
