package tracker

import "github.com/huaweicloud/golangsdk"

const rootPath = "tracker"
const trackerName = "system"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, trackerName)
}
