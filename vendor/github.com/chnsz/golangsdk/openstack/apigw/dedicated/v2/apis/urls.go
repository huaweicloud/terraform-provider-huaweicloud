package apis

import (
	"github.com/chnsz/golangsdk"
)

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "apis")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(rootPath, instanceId, "apis", apiId)
}

func releaseURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "apis", "action")
}

func batchPublishURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "apis", "publish")
}

func publishVersionURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(rootPath, instanceId, "apis", "publish", apiId)
}

func showHistoryDetailURL(c *golangsdk.ServiceClient, instanceId, versionId string) string {
	return c.ServiceURL(rootPath, instanceId, "apis", "versions", versionId)
}
