package runtimestacks

import "github.com/chnsz/golangsdk"

const rootPath = "cas/runtimestacks"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}
