package groups

import (
	"github.com/chnsz/golangsdk"
)

// endpoint/queues/{queue_id}/groups
const resourcePathQueues = "queues"
const resourcePathGroups = "groups"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient, queueID string) string {
	return client.ServiceURL(client.ProjectID, resourcePathQueues, queueID, resourcePathGroups)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, queueID string, groupID string) string {
	return client.ServiceURL(client.ProjectID, resourcePathQueues, queueID, resourcePathGroups, groupID)
}

// listURL will build the list url of list function
func listURL(client *golangsdk.ServiceClient, queueID string) string {
	return client.ServiceURL(client.ProjectID, resourcePathQueues, queueID, "resourcePathGroups")
}
