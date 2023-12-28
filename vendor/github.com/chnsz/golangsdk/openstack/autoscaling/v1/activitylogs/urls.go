package activitylogs

import (
	"github.com/chnsz/golangsdk"
)

func listURL(c *golangsdk.ServiceClient, groupID string) string {
	return c.ServiceURL("scaling_activity_log", groupID)
}
