package users

import "github.com/chnsz/golangsdk"

const rootPath = "OS-USER"

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "users")
}

func updateURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL(rootPath, "users", userID)
}

func getURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL(rootPath, "users", userID)
}
