package queues

import "github.com/chnsz/golangsdk"

func propertyURL(c *golangsdk.ServiceClient, queueName string) string {
	return c.ServiceURL("queues", queueName, "properties")
}
