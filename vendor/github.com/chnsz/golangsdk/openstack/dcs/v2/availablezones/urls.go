package availablezones

import (
	"github.com/chnsz/golangsdk"
)

const resourcePath = "available-zones"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}
