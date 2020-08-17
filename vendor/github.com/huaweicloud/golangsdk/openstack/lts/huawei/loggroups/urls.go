package loggroups

import "github.com/huaweicloud/golangsdk"

const rootPath = "groups"

// createURL will build the url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}

// updateURL will build the url of updation
func updateURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(rootPath, id)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(rootPath, id)
}

// listURL will build the get url of get list function
func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}
