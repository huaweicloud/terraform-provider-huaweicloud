package services

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("workspaces")
}

func authConfigURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("assist-auth-config/method-config")
}

func lockStatusURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("workspaces/lock-status")
}
