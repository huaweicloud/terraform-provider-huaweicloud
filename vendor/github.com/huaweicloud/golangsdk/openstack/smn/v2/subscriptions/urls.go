package subscriptions

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient, topicUrn string) string {
	return c.ServiceURL("topics", topicUrn, "subscriptions")
}

func deleteURL(c *golangsdk.ServiceClient, subscriptionUrn string) string {
	return c.ServiceURL("subscriptions", subscriptionUrn)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("subscriptions?offset=0&limit=100")
}

func listFromTopicURL(c *golangsdk.ServiceClient, topicUrn string) string {
	return c.ServiceURL("topics", topicUrn, "subscriptions?offset=0&limit=100")
}
