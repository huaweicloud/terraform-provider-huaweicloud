package sms

import(
	"github.com/chnsz/golangsdk"
)

const resourcePath = "sms"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}