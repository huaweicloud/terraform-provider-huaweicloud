package accesspolicies

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("access-policy")
}

func resourceURL(c *golangsdk.ServiceClient, policyId string) string {
	return c.ServiceURL("access-policy", policyId, "objects")
}
