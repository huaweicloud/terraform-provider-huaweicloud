package apis

import (
	"fmt"

	"github.com/chnsz/golangsdk"
)

const rootPath = "instances"

func buildRootPath(instanceId string) string {
	return fmt.Sprintf("instances/%s/apis", instanceId)
}

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId))
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(buildRootPath(instanceId), apiId)
}

func releaseURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "action")
}

func batchPublishURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "publish")
}

func publishVersionURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "publish", apiId)
}

func showHistoryDetailURL(c *golangsdk.ServiceClient, instanceId, versionId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "version", versionId)
}
