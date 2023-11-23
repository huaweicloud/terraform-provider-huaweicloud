package ipaddress

import "github.com/chnsz/golangsdk"

func baseUrl(c *golangsdk.ServiceClient, endpointID string) string {
	return c.ServiceURL("endpoints", endpointID, "ipaddresses")
}

func resourceUrl(c *golangsdk.ServiceClient, endpointID string, ipaddress string) string {
	return c.ServiceURL("endpoints", endpointID, "ipaddresses", ipaddress)
}
