package dependencies

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("fgs", "dependencies")
}

func resourceURL(client *golangsdk.ServiceClient, dependId string) string {
	return client.ServiceURL("fgs", "dependencies", dependId)
}
