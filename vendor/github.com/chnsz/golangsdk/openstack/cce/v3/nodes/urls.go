package nodes

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "clusters"
	resourcePath = "nodes"
)

func rootURL(c *golangsdk.ServiceClient, clusterid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, clusterid, nodeid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath, nodeid)
}

func getJobURL(c *golangsdk.ServiceClient, jobid string) string {
	return c.ServiceURL("jobs", jobid)
}

func removeNodeURL(c *golangsdk.ServiceClient, clusterid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath, "operation/remove")
}

func addNodeURL(c *golangsdk.ServiceClient, clusterid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath, "add")
}

func resetNodeURL(c *golangsdk.ServiceClient, clusterid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath, "reset")
}
