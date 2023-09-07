package custom

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("sources")
}

func resourceURL(c *golangsdk.ServiceClient, sourceId string) string {
	return c.ServiceURL("sources", sourceId)
}
