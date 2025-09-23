package configurations

import "github.com/chnsz/golangsdk"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("configurations")
}

func applyURL(c *golangsdk.ServiceClient, configurationId string) string {
	return c.ServiceURL("configurations", configurationId, "apply")
}
