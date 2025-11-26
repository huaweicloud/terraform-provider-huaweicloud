package instances

import "github.com/chnsz/golangsdk"

// endpoint/instances
const resourcePath = "instances"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}

func createURLWithEngine(engine string, client *golangsdk.ServiceClient) string {
	return client.ServiceURL(engine, client.ProjectID, resourcePath)
}

func createInstanceURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, "kafka", resourcePath)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

// updateURL will build the url of update
func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}

func extend(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id, "extend")
}

func extendInstanceURL(client *golangsdk.ServiceClient, instanceId string) string {
	return client.ServiceURL(client.ProjectID, "kafka", resourcePath, instanceId, "extend")
}

func crossVpcURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id, "crossvpc/modify")
}

// autoTopicURL will build the url of UpdateAutoTopic
func autoTopicURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "autotopic")
}

// resetPasswordURL will build the url of resetting password
func resetPasswordURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "password")
}

func configurationsURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id, "configs")
}

func actionURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "instances", "action")
}

func tasksURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id, "tasks")
}

func taskURL(c *golangsdk.ServiceClient, instanceID, taskID string) string {
	return c.ServiceURL(c.ProjectID, "instances", instanceID, "tasks", taskID)
}
