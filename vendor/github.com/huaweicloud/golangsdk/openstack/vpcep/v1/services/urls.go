package services

import "github.com/huaweicloud/golangsdk"

const (
	rootPath = "vpc-endpoint-services"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, serviceID string) string {
	return c.ServiceURL(rootPath, serviceID)
}
