package backups

import "github.com/chnsz/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "backups/policy/update")
}

func updateEncryptionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "backups/encryption")
}

func getEncryptionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "backups/encryption")
}
