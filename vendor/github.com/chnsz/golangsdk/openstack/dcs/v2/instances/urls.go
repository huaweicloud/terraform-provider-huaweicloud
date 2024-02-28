package instances

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "instances")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id)
}

func resizeResourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id, "resize")
}

func resizePrePaidResourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "orders/instances", id, "resize")
}

func updatePasswordURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id, "password")
}

func restartOrFlushInstanceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "instances/status")
}

func configurationsURL(c *golangsdk.ServiceClient, instancesId string) string {
	return c.ServiceURL(c.ProjectID, "instances", instancesId, "configs")
}

func sslURL(client *golangsdk.ServiceClient, instancesId string) string {
	return client.ServiceURL(client.ProjectID, "instances", instancesId, "ssl")
}
