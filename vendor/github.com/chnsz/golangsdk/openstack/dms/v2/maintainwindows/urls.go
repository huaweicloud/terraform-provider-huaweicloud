package maintainwindows

import (
	"github.com/chnsz/golangsdk"
)

// endpoint/instances/maintain-windows
const resourcePath = "instances/maintain-windows"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}
