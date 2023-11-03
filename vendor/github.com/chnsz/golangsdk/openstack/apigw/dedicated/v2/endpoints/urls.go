package endpoints

import "github.com/chnsz/golangsdk"

func listURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "vpc-endpoint/permissions")
}

func addURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "vpc-endpoint/permissions/batch-add")
}

func deleteURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "vpc-endpoint/permissions/batch-delete")
}
