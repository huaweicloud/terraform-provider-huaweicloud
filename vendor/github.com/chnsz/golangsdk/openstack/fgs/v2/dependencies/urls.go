package dependencies

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("fgs", "dependencies")
}

func resourceURL(client *golangsdk.ServiceClient, dependId string) string {
	return client.ServiceURL("fgs", "dependencies", dependId)
}

func rootVersionURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("fgs/dependencies/version")
}

func resourceVersionURL(client *golangsdk.ServiceClient, dependId, version string) string {
	return client.ServiceURL("fgs/dependencies", dependId, "version", version)
}

func resourceVersionsURL(client *golangsdk.ServiceClient, dependId string) string {
	return client.ServiceURL("fgs/dependencies", dependId, "version")
}
