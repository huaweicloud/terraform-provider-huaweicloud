package logstreams

import "github.com/huaweicloud/golangsdk"

const (
	resourcePath = "streams"
	rootPath     = "groups"
)

// createURL will build the url of creation
func createURL(client *golangsdk.ServiceClient, groupId string) string {
	return client.ServiceURL(rootPath, groupId, resourcePath)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, groupId string, id string) string {
	return client.ServiceURL(rootPath, groupId, resourcePath, id)
}

// listURL will build the get url of get function
func listURL(client *golangsdk.ServiceClient, groupId string) string {
	return client.ServiceURL(rootPath, groupId, resourcePath)
}
