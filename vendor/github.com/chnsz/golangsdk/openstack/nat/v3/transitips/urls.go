package transitips

import "github.com/chnsz/golangsdk"

const rootPath = "private-nat/transit-ips"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, transitIpId string) string {
	return c.ServiceURL(rootPath, transitIpId)
}
