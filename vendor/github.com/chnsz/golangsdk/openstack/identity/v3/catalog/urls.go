package catalog

import "github.com/chnsz/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("auth/catalog")
}
