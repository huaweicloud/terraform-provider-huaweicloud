package applications

import "github.com/chnsz/golangsdk"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "apps")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, appId string) string {
	return c.ServiceURL(rootPath, instanceId, "apps", appId)
}

func resetSecretURL(c *golangsdk.ServiceClient, instanceId, appId string) string {
	return c.ServiceURL(rootPath, instanceId, "apps/secret", appId)
}

func codeURL(c *golangsdk.ServiceClient, instanceId, appId string) string {
	return c.ServiceURL(rootPath, instanceId, "apps", appId, "app-codes")
}

func codeResourceURL(c *golangsdk.ServiceClient, instanceId, appId, codeId string) string {
	return c.ServiceURL(rootPath, instanceId, "apps", appId, "app-codes", codeId)
}
