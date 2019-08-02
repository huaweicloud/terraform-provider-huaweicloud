package volumeactions

import "github.com/huaweicloud/golangsdk"

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("volumes", id, "action")
}
