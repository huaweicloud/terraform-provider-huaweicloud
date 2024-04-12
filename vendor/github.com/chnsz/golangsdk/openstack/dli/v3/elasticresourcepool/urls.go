package elasticresourcepool

import "github.com/chnsz/golangsdk"

func associateQueueURl(c *golangsdk.ServiceClient, name string) string {
	return c.ServiceURL("elastic-resource-pools", name, "queues")
}

func queueScalingPolicyURL(c *golangsdk.ServiceClient, elasticResourcePoolName, queueName string) string {
	return c.ServiceURL("elastic-resource-pools", elasticResourcePoolName, "queues", queueName)
}
