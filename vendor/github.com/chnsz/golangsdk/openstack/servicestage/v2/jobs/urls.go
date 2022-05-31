package jobs

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, jobId string) string {
	return client.ServiceURL("cas/jobs", jobId)
}
