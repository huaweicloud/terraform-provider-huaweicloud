package queues

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

// endpoint/queues
const resourcePath = "queues"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, id string, includeDeadLetter bool) string {
	return client.ServiceURL(resourcePath, fmt.Sprintf("%s?include_deadletter=%t", id, includeDeadLetter))
}

// listURL will build the list url of list function
func listURL(client *golangsdk.ServiceClient, includeDeadLetter bool) string {
	return client.ServiceURL(fmt.Sprintf("%s?include_deadletter=%t", resourcePath, includeDeadLetter))
}
