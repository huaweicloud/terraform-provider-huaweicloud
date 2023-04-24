package signs

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "signs")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, signatureId string) string {
	return c.ServiceURL("instances", instanceId, "signs", signatureId)
}

func bindURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "sign-bindings")
}

func listBindURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "sign-bindings/binded-apis")
}

func unbindURL(c *golangsdk.ServiceClient, instanceId, bindId string) string {
	return c.ServiceURL("instances", instanceId, "sign-bindings", bindId)
}
