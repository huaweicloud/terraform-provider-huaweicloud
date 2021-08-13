package apis

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
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
