package customer_gateways

import "github.com/chnsz/golangsdk"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("customer-gateways")
}
