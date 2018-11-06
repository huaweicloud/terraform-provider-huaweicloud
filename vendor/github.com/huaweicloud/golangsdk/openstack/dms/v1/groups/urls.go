package groups

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

// endpoint/queues/{queue_id}/groups
const resourcePathQueues = "queues"
const resourcePathGroups = "groups"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient, queueID string) string {
	return client.ServiceURL(resourcePathQueues, queueID, resourcePathGroups)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, queueID string, groupID string) string {
	return client.ServiceURL(resourcePathQueues, queueID, resourcePathGroups, groupID)
}

// listURL will build the list url of list function
func listURL(client *golangsdk.ServiceClient, queueID string, includeDeadLetter bool) string {
	return client.ServiceURL(resourcePathQueues, queueID, fmt.Sprintf("%s?include_deadletter=%t", resourcePathGroups, includeDeadLetter))
}
