package environments

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("cas", "environments")
}

func detailURL(client *golangsdk.ServiceClient, envId string) string {
	return client.ServiceURL("cas", "environments", envId)
}

func resourceURL(client *golangsdk.ServiceClient, envId string) string {
	return client.ServiceURL("cas", "environments", envId, "resources")
}
