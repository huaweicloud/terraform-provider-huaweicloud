package job

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("jobs/submit-job")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("job-executions", id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("job-exes", id)
}
