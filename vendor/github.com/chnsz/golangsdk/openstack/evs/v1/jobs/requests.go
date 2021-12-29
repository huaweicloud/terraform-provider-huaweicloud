package jobs

import (
	"github.com/chnsz/golangsdk"
)

// GetJobDetails retrieves the Job with the provided jobID. To extract the Job object
// from the response, call the ExtractJob method on the GetResult.
func GetJobDetails(client *golangsdk.ServiceClient, jobID string) (r GetResult) {
	_, r.Err = client.Get(jobURL(client, jobID), &r.Body, nil)
	return
}
