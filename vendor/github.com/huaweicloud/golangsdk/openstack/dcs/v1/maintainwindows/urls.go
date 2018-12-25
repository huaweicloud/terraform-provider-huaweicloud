package maintainwindows

import (
	"strings"

	"github.com/huaweicloud/golangsdk"
)

// endpoint/instances/maintain-windows
const resourcePath = "instances/maintain-windows"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	// remove projectid from endpoint
	return strings.Replace(client.ServiceURL(resourcePath), "/"+client.ProjectID, "", -1)
}
