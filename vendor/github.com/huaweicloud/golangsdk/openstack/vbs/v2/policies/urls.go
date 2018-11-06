package policies

import "github.com/huaweicloud/golangsdk"

const (
	backupRootPath     = "backuppolicy"
	policyResourcePath = "backuppolicyresources"
)

func commonURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(backupRootPath)
}

func resourceURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL(backupRootPath, policyID)
}

func associateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(policyResourcePath)
}

func disassociateURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL(policyResourcePath, policyID, "deleted_resources")
}
