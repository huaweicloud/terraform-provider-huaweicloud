package topics

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("topics")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("topics", id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("topics")
}

func getPoliciesURL(c *golangsdk.ServiceClient, id, policyName string) string {
	return c.ServiceURL("topics", id, "attributes?name="+policyName)
}

func updatePoliciesURL(c *golangsdk.ServiceClient, id, policyName string) string {
	return c.ServiceURL("topics", id, "attributes", policyName)
}
