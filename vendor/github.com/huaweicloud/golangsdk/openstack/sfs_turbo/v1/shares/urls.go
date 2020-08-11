package shares

import "github.com/huaweicloud/golangsdk"

// createURL used to assemble the URI of creating API
func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("sfs-turbo/shares")
}

// resourceURL used to assemble the URI of deleting API or querying details API
func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("sfs-turbo/shares", id)
}

// listURL used to assemble the URI of querying all file system details API
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("sfs-turbo/shares", "detail")
}

// For mamage the specified file system, e.g.: extend
func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("sfs-turbo/shares", id, "action")
}
