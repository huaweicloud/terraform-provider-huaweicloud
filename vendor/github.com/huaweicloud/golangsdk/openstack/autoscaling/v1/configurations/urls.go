package configurations

import (
	"github.com/huaweicloud/golangsdk"
)

const resourcePath = "scaling_configuration"

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}
