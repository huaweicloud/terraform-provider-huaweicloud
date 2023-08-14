package jobs

import "github.com/chnsz/golangsdk"

func jobURL(c *golangsdk.ServiceClient, jobID string) string {
	return c.ServiceURL("jobs", jobID)
}
