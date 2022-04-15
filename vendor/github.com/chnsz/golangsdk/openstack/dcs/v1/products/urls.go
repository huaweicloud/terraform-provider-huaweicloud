package products

import (
	"github.com/chnsz/golangsdk"
)

// endpoint/products
const resourcePath = "products"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}
