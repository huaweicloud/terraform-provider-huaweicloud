package endpoints

import "github.com/huaweicloud/golangsdk"

const (
	rootPath = "vpc-endpoints"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, endpointID string) string {
	return c.ServiceURL(rootPath, endpointID)
}
