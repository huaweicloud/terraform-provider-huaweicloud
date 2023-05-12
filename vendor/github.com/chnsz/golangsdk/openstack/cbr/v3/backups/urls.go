package backups

import "github.com/chnsz/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, backupId string) string {
	return c.ServiceURL("backups", backupId)
}
