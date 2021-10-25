package oidcconfig

import "github.com/chnsz/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, idpID string) string {
	return c.ServiceURL("v3.0", "OS-FEDERATION", "identity-providers", idpID, "openid-connect-config")
}
