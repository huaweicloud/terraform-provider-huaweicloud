package subscriptions

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("subscriptions")
}

func resourceURL(c *golangsdk.ServiceClient, subscriptionId string) string {
	return c.ServiceURL("subscriptions", subscriptionId)
}
