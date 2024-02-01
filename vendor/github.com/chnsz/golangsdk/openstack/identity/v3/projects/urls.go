package projects

import "github.com/chnsz/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("projects")
}

func getURL(client *golangsdk.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("projects")
}

func deleteURL(client *golangsdk.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func updateURL(client *golangsdk.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func updateStatusURL(c *golangsdk.ServiceClient, projectID string) string {
	return c.ServiceURL("projects", projectID)
}
