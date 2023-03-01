package gateways

import "github.com/chnsz/golangsdk"

const rootPath = "nat_gateways"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, gatewayId string) string {
	return c.ServiceURL(rootPath, gatewayId)
}
