package signs

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "signs")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, signatureId string) string {
	return c.ServiceURL("instances", instanceId, "signs", signatureId)
}
