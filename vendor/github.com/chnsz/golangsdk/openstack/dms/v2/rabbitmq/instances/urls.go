package instances

import "github.com/chnsz/golangsdk"

// endpoint/instances
const resourcePath = "instances"

const rabbitMqEngine = "rabbitmq"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}

// createWithEngineURL will build the rest query url of creation
func createWithEngineURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rabbitMqEngine, client.ProjectID, resourcePath)
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
	return client.ServiceURL(rabbitMqEngine, client.ProjectID, resourcePath, id, "extend")
}

// resetPasswordURL will build the url of resetting password
func resetPasswordURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "password")
}
