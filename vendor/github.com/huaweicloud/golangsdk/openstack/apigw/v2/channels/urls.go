package channels

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

func buildRootPath(instanceId string) string {
	return fmt.Sprintf("instances/%s/vpc-channels", instanceId)
}

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId))
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, chanId string) string {
	return c.ServiceURL(buildRootPath(instanceId), chanId)
}

func membersURL(c *golangsdk.ServiceClient, instanceId, chanId string) string {
	return c.ServiceURL(buildRootPath(instanceId), chanId, "members")
}

func memberURL(c *golangsdk.ServiceClient, instanceId, chanId, memberId string) string {
	return c.ServiceURL(buildRootPath(instanceId), chanId, "members", memberId)
}
