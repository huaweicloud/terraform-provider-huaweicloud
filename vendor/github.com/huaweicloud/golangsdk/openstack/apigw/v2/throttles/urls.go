package throttles

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

func buildRootPath(instanceId string) string {
	return fmt.Sprintf("instances/%s/throttles", instanceId)
}

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId))
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, policyId string) string {
	return c.ServiceURL(buildRootPath(instanceId), policyId)
}

func specRootURL(c *golangsdk.ServiceClient, instanceId, policyId string) string {
	return c.ServiceURL(buildRootPath(instanceId), policyId, "throttle-specials")
}

func specResourceURL(c *golangsdk.ServiceClient, instanceId, policyId, strategyId string) string {
	return c.ServiceURL(buildRootPath(instanceId), policyId, "throttle-specials", strategyId)
}
