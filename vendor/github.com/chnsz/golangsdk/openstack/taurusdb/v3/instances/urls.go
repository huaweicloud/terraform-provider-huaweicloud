package instances

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID)
}

func getURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func enlargeURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "nodes/enlarge")
}

func deleteReplicaURL(c *golangsdk.ServiceClient, instanceID, nodeID string) string {
	return c.ServiceURL("instances", instanceID, "nodes", nodeID)
}

func nameURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "name")
}

func passwordURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "password")
}

func volumeURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "volume/extend")
}

func proxyURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "proxy")
}

func proxyEnlargeURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "proxy/enlarge")
}

func actionURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "action")
}

func jobURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("jobs")
}

func listDehURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("dedicated-resources")
}
