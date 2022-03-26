package availablezones

import (
	"github.com/chnsz/golangsdk"
)

// endpoint/availablezones
const resourcePath = "availableZones"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}
