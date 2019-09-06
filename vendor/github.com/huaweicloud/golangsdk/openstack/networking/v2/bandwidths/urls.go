package bandwidths

import "github.com/huaweicloud/golangsdk"

func PostURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "bandwidths")
}

func BatchPostURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "batch-bandwidths")
}
func UpdateURL(c *golangsdk.ServiceClient, ID string) string {
	return c.ServiceURL(c.ProjectID, "bandwidths", ID)
}

func DeleteURL(c *golangsdk.ServiceClient, ID string) string {
	return c.ServiceURL(c.ProjectID, "bandwidths", ID)
}

func InsertURL(c *golangsdk.ServiceClient, ID string) string {
	return c.ServiceURL(c.ProjectID, "bandwidths", ID, "insert")
}

func RemoveURL(c *golangsdk.ServiceClient, ID string) string {
	return c.ServiceURL(c.ProjectID, "bandwidths", ID, "remove")
}
