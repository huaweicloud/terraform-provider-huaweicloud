package queues

import "github.com/huaweicloud/golangsdk"

const (
	resourcePath = "queues"
)

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, queueName string) string {
	return c.ServiceURL(resourcePath, queueName)
}

func queryAllURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
