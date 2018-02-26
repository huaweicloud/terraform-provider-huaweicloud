package backendmember

import "github.com/huaweicloud/golangsdk"

const (
	rootPath     = "elbaas"
	resourcePath = "listeners"
)

func addURL(c *golangsdk.ServiceClient, listener_id string) string {
	return c.ServiceURL(rootPath, resourcePath, listener_id, "members")
}

func removeURL(c *golangsdk.ServiceClient, listener_id string) string {
	return c.ServiceURL(rootPath, resourcePath, listener_id, "members", "action")
}

func resourceURL(c *golangsdk.ServiceClient, listener_id string, id string) string {
	return c.ServiceURL(rootPath, resourcePath, listener_id, "members?id="+id)
}
