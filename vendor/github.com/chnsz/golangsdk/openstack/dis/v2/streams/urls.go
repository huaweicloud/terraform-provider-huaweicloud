package streams

import "github.com/chnsz/golangsdk"

const (
	resourcePath = "streams"
	policiesPath = "policies"
)

// createURL POST /v2/{project_id}/streams
// ListURL   GET  /v2/{project_id}/streams
func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

// GetURL GET /v2/{project_id}/streams/{stream_name}
// deleteURL DELETE /v2/{project_id}/streams/{stream_name}
// ChangePartitionQuantityURL PUT /v2/{project_id}/streams/{stream_name}
func resourceURL(c *golangsdk.ServiceClient, streamName string) string {
	return c.ServiceURL(resourcePath, streamName)
}

// policiesURL POST /v2/{project_id}/streams/{stream_name}/policies
// policiesURL GET /v2/{project_id}/streams/{stream_name}/policies
func policiesURL(c *golangsdk.ServiceClient, streamName string) string {
	return c.ServiceURL(resourcePath, streamName, policiesPath)
}
