package flowlogs

import "github.com/chnsz/golangsdk"

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "fl/flow_logs")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "fl/flow_logs")
}

func GetURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "fl/flow_logs", id)
}

func UpdateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "fl/flow_logs", id)
}

func DeleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "fl/flow_logs", id)
}
