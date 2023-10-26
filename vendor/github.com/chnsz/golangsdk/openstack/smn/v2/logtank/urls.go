package logtank

import (
	"github.com/chnsz/golangsdk"
)

// baseURL be used as CreateLogtank or ListLogtank request url
func baseURL(c *golangsdk.ServiceClient, topicUrn string) string {
	return c.ServiceURL("topics", topicUrn, "logtanks")
}

// resourceURL be used as UpdateLogtank or DeleteLogtank request url
func resourceURL(c *golangsdk.ServiceClient, topicUrn string, logTankID string) string {
	return c.ServiceURL("topics", topicUrn, "logtanks", logTankID)
}
