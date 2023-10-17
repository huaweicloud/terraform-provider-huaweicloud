package checkpoints

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("checkpoints")
}

func resourceURL(c *golangsdk.ServiceClient, checkpointId string) string {
	return c.ServiceURL("checkpoints", checkpointId)
}
