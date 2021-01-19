package backups

import "github.com/huaweicloud/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "backups/policy")
}
