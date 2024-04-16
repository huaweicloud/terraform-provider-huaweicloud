package flinkjob

import "github.com/chnsz/golangsdk"

func streamGraphURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL("streaming/jobs", jobId, "gen-graph")
}
