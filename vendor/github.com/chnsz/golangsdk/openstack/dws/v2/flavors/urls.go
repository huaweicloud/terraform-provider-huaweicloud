package flavors

import "github.com/chnsz/golangsdk"

// listNodeTypesURL /v2/{project_id}/node-types
func listNodeTypesURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("node-types")
}
