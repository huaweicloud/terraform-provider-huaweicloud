package volumeactions

import "github.com/chnsz/golangsdk"

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("volumes", id, "action")
}
