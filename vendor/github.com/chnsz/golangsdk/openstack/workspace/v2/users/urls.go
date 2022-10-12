package users

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("users")
}

func resourceURL(c *golangsdk.ServiceClient, userId string) string {
	return c.ServiceURL("users", userId)
}
