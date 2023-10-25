package resourcetags

import "github.com/chnsz/golangsdk"

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("resource-tags/batch-create")
}

func queryURL(client *golangsdk.ServiceClient, resourceId string) string {
	return client.ServiceURL("resources", resourceId, "tags")
}

func deleteURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("resource-tags/batch-delete")
}
