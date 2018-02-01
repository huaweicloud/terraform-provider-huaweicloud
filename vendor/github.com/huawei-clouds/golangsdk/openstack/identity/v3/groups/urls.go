package groups

import "github.com/huawei-clouds/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("groups")
}

func getURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("groups")
}

func updateURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID)
}

func deleteURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID)
}
