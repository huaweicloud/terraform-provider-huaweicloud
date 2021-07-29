package jobs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

func buildRootPath(clusterId string) string {
	return fmt.Sprintf("clusters/%s/job-executions", clusterId)
}

func rootURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL(buildRootPath(clusterId))
}

func resourceURL(c *golangsdk.ServiceClient, clusterId, jobId string) string {
	return c.ServiceURL(buildRootPath(clusterId), jobId)
}

func deleteURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL(buildRootPath(clusterId), "batch-delete")
}
