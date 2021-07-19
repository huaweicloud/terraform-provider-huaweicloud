package responses

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

func buildResponsesPath(instanceId, groupId string) string {
	return fmt.Sprintf("instances/%s/api-groups/%s/gateway-responses", instanceId, groupId)
}

func rootURL(c *golangsdk.ServiceClient, path string) string {
	return c.ServiceURL(path)
}

func resourceURL(c *golangsdk.ServiceClient, path, respId string) string {
	return c.ServiceURL(path, respId)
}

func specResponsesURL(c *golangsdk.ServiceClient, path, respId, respType string) string {
	return c.ServiceURL(path, respId, respType)
}
