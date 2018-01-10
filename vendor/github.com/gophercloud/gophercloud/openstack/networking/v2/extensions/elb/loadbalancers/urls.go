package loadbalancers

import "github.com/gophercloud/gophercloud"

const (
	rootPath     = "elbaas"
	resourcePath = "loadbalancers"
)

func rootURL(c *gophercloud.ServiceClient1) string {
	return c.ServiceURL(c.ProjectID, rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient1, id string) string {
	return c.ServiceURL(c.ProjectID, rootPath, resourcePath, id)
}
