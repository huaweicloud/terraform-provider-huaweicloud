package protocols

import "github.com/chnsz/golangsdk"

func root(c *golangsdk.ServiceClient, idpID string) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "identity_providers", idpID, "protocols")
}

func resourceURL(c *golangsdk.ServiceClient, idpID string, protocolID string) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "identity_providers", idpID, "protocols", protocolID)
}
