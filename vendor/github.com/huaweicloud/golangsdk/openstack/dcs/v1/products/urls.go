package products

import (
	"strings"

	"github.com/huaweicloud/golangsdk"
)

// endpoint/products
const resourcePath = "products"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	// remove projectid from endpoint
	return strings.Replace(client.ServiceURL(resourcePath), "/"+client.ProjectID, "", -1)
}
