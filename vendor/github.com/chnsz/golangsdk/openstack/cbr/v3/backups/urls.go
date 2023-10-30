package backups

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("backups")
}

func resourceURL(c *golangsdk.ServiceClient, backupId string) string {
	return c.ServiceURL("backups", backupId)
}
