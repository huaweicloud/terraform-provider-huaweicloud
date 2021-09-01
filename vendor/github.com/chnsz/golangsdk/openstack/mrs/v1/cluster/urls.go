package cluster

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("run-job-flow")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("cluster_infos", id)
}

// listHostsURL /v1.1/{project_id}/clusters/{cluster_id}/hosts
func listHostsURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "hosts")
}
