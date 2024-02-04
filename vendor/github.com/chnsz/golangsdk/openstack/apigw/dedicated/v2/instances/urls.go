package instances

import "github.com/chnsz/golangsdk"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}

func egressURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "nat-eip")
}

func ingressURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "eip")
}

func elbIngressURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "ingress-eip")
}

func featureURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "features")
}

func modifyTagsURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "instance-tags/action")
}

func queryTagsURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "instance-tags")
}
