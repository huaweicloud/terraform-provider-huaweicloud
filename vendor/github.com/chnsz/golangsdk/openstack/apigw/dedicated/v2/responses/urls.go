package responses

import (
	"github.com/chnsz/golangsdk"
)

func rootURL(c *golangsdk.ServiceClient, instanceId, groupId string) string {
	return c.ServiceURL("instances", instanceId, "api-groups", groupId, "gateway-responses")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, groupId, respId string) string {
	return c.ServiceURL("instances", instanceId, "api-groups", groupId, "gateway-responses", respId)
}

func specResponsesURL(c *golangsdk.ServiceClient, instanceId, groupId, respId, respType string) string {
	return c.ServiceURL("instances", instanceId, "api-groups", groupId, "gateway-responses", respId, respType)
}
