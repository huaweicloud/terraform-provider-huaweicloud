package jobs

import "github.com/chnsz/golangsdk"

const rootPath = "cas/jobs"

func rootURL(client *golangsdk.ServiceClient, jobId string) string {
	return client.ServiceURL(rootPath, jobId)
}
