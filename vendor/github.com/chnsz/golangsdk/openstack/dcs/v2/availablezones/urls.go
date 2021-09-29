package availablezones

import (
	"strings"

	"github.com/chnsz/golangsdk"
)

const resourcePath = "available-zones"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	// remove projectId from endpoint
	return strings.Replace(client.ServiceURL(resourcePath), "/"+client.ProjectID, "", -1)
}
