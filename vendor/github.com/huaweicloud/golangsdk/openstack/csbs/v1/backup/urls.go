package backup

import "github.com/huaweicloud/golangsdk"

const rootPath = "providers"
const ProviderID = "fc4d5750-22e7-4798-8a46-f48f62c4c1da"

func rootURL(c *golangsdk.ServiceClient, resourceid string) string {
	return c.ServiceURL(rootPath, ProviderID, "resources", resourceid, "action")
}

func resourceURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, ProviderID, "resources", "action")
}

func getURL(c *golangsdk.ServiceClient, checkpoint_item_id string) string {
	return c.ServiceURL("checkpoint_items", checkpoint_item_id)
}

func deleteURL(c *golangsdk.ServiceClient, checkpoint_id string) string {
	return c.ServiceURL(rootPath, ProviderID, "checkpoints", checkpoint_id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("checkpoint_items")
}
