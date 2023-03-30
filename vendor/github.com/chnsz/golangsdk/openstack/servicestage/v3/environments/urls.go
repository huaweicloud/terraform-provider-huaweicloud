package environments

import "github.com/chnsz/golangsdk"

const rootPath = "cas/environments"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}

func detailURL(client *golangsdk.ServiceClient, envId string) string {
	return client.ServiceURL(rootPath, envId)
}

func resourceURL(client *golangsdk.ServiceClient, envId string) string {
	return client.ServiceURL(rootPath, envId, "resources")
}
