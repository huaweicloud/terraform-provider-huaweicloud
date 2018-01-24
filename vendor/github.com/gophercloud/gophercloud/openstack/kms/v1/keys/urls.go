package keys

import "github.com/gophercloud/gophercloud"

func getURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("describe-key")
}

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("create-key")
}

func deleteURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("schedule-key-deletion")
}

func updateAliasURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("update-key-alias")
}

func updateDesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("update-key-description")
}

func dataEncryptURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("create-datakey")
}

func encryptDEKURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("encrypt-datakey")
}

func enableKeyURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("enable-key")
}

func disableKeyURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("disable-key")
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("list-keys")
}
