package users

import "github.com/huaweicloud/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("users")
}

func getURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("users")
}

func updateURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID)
}

func deleteURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID)
}

func listGroupsURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID, "groups")
}

func listProjectsURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID, "projects")
}

func listInGroupURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID, "users")
}

func membershipURL(client *golangsdk.ServiceClient, groupID string, userID string) string {
	return client.ServiceURL("groups", groupID, "users", userID)
}
