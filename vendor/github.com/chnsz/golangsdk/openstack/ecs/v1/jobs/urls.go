package jobs

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL("jobs", jobId)
}
