package certificates

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "elb"
	resourcePath = "certificates"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}
