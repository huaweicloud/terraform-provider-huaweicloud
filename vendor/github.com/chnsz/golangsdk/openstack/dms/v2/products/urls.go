package products

import (
	"github.com/chnsz/golangsdk"
)

// endpoint/products
const resourcePath = "products"

func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

func listURL(client *golangsdk.ServiceClient, engineType string) string {
	return client.ServiceURL(engineType, resourcePath)
}
