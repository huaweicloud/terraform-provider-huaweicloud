package keypairs

import "github.com/chnsz/golangsdk"

const (
	resourcePath = "keypairs"
)

// associateURL /v3/{project_id}/keypairs/associate
func associateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "associate")
}

// disassociateURL /v3/{project_id}/keypairs/disassociate
func disassociateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "disassociate")
}

func getTaskURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("tasks", id)
}
