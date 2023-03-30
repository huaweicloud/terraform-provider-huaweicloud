package components

import (
	"fmt"
	"github.com/chnsz/golangsdk"
)

const rootPath = "cas/applications/%s/components"
const componentPath = "cas/components"

func rootURL(c *golangsdk.ServiceClient, appId string) string {
	return c.ServiceURL(fmt.Sprintf(rootPath, appId))
}

func resourceURL(c *golangsdk.ServiceClient, appId, componentId string) string {
	return c.ServiceURL(fmt.Sprintf(rootPath, appId), componentId)
}

func componentURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(componentPath)
}

func actionURL(c *golangsdk.ServiceClient, appId, componentId string) string {
	return c.ServiceURL(fmt.Sprintf(rootPath, appId), componentId, "action")
}

func recordURL(c *golangsdk.ServiceClient, appId, componentId string) string {
	return c.ServiceURL(fmt.Sprintf(rootPath, appId), componentId, "records")
}
