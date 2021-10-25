package providers

import "github.com/chnsz/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "identity_providers", id)
}
