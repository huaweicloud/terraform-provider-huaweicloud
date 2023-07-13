package certificates

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("certificates")
}

func resourceURL(c *golangsdk.ServiceClient, certificateId string) string {
	return c.ServiceURL("certificates", certificateId)
}
