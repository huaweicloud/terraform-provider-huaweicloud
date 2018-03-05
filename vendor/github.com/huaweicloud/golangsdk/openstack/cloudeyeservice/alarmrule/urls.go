package alarmrule

import "github.com/huaweicloud/golangsdk"

const (
	rootPath = "alarms"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, id)
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, id, "action")
}
