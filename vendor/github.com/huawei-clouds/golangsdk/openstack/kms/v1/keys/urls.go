package keys

import "github.com/huawei-clouds/golangsdk"

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("describe-key")
}

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("create-key")
}

func deleteURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("schedule-key-deletion")
}

func updateAliasURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("update-key-alias")
}

func updateDesURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("update-key-description")
}

func dataEncryptURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("create-datakey")
}

func encryptDEKURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("encrypt-datakey")
}

func enableKeyURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("enable-key")
}

func disableKeyURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("disable-key")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("list-keys")
}
