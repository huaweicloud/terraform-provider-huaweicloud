package extensions

import "github.com/huaweicloud/golangsdk"

// ExtensionURL generates the URL for an extension resource by name.
func ExtensionURL(c *golangsdk.ServiceClient, name string) string {
	return c.ServiceURL("extensions", name)
}

// ListExtensionURL generates the URL for the extensions resource collection.
func ListExtensionURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("extensions")
}
