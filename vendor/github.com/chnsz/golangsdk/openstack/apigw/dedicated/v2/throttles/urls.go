package throttles

import (
	"github.com/chnsz/golangsdk"
)

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "throttles")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, policyId string) string {
	return c.ServiceURL("instances", instanceId, "throttles", policyId)
}

func specRootURL(c *golangsdk.ServiceClient, instanceId, policyId string) string {
	return c.ServiceURL("instances", instanceId, "throttles", policyId, "throttle-specials")
}

func specResourceURL(c *golangsdk.ServiceClient, instanceId, policyId, strategyId string) string {
	return c.ServiceURL("instances", instanceId, "throttles", policyId, "throttle-specials", strategyId)
}

func bindURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "throttle-bindings")
}

func listBindURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "throttle-bindings", "binded-apis")
}

func unbindURL(c *golangsdk.ServiceClient, instanceId, bindId string) string {
	return c.ServiceURL("instances", instanceId, "throttle-bindings", bindId)
}

func batchUnbindURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "throttle-bindings")
}
