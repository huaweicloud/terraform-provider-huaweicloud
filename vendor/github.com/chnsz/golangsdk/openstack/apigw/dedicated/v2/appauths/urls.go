package appauths

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "app-auths")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, authId string) string {
	return c.ServiceURL("instances", instanceId, "app-auths", authId)
}

func listAuthorizedURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "app-auths/binded-apis")
}

func listUnathorizedURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "app-auths/unbinded-apis")
}
