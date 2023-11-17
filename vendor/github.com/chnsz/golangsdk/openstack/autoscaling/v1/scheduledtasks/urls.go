package scheduledtasks

import (
	"github.com/chnsz/golangsdk"
)

const resourcePath = "scaling-groups"

func rootURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL(resourcePath, groupID, "scheduled-tasks")
}

func resourceURL(client *golangsdk.ServiceClient, groupID, taskID string) string {
	return client.ServiceURL(resourcePath, groupID, "scheduled-tasks", taskID)
}
