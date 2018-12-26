package cluster

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("run-job-flow")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("cluster_infos", id)
}
