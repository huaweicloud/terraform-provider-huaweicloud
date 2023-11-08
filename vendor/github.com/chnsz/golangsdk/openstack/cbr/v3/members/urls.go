package members

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, backupId string) string {
	return c.ServiceURL("backups", backupId, "members")
}

func resourceURL(c *golangsdk.ServiceClient, backupId, memberId string) string {
	return c.ServiceURL("backups", backupId, "members", memberId)
}
