package alarmrule

import "github.com/chnsz/golangsdk"

const (
	rootPath = "alarms"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func listURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("alarms?alarm_id=" + id)
}

func actionURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "action")
}

func batchResourcesURL(c *golangsdk.ServiceClient, id, operation string) string {
	return c.ServiceURL(rootPath, id, "resources", operation)
}

func policiesURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "policies")
}

func resourcesURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "resources")
}

func deleteURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "batch-delete")
}
