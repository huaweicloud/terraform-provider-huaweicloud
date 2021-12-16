package products

import (
	"strings"

	"github.com/chnsz/golangsdk"
)

// endpoint/products
const resourcePath = "products"

func getURL(client *golangsdk.ServiceClient) string {
	// remove projectid from endpoint
	return strings.Replace(client.ServiceURL(resourcePath), "/"+client.ProjectID, "", -1)
}
