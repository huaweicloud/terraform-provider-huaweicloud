package topics

import (
	"github.com/chnsz/golangsdk"
)

const (
	resourcePath = "instances"
	topicPath    = "topics"
)

// rootURL will build the url of create, update and list
func rootURL(client *golangsdk.ServiceClient, instanceID string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, instanceID, topicPath)
}

// getURL will build the url of get
func getURL(client *golangsdk.ServiceClient, instanceID, topic string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, instanceID, "management", topicPath, topic)
}

// deleteURL will build the url of delete
func deleteURL(client *golangsdk.ServiceClient, instanceID string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, instanceID, topicPath, "delete")
}
