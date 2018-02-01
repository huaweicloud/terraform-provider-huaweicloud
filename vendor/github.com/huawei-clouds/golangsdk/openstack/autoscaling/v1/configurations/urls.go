package configurations

import (
	"github.com/huawei-clouds/golangsdk"
)

const resourcePath = "scaling_configuration"

func createURL(client *golangsdk.ServiceClientExtension) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}

func getURL(client *golangsdk.ServiceClientExtension, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

func deleteURL(client *golangsdk.ServiceClientExtension, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

func listURL(client *golangsdk.ServiceClientExtension) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}
