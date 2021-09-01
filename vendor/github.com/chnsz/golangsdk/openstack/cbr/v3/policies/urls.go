package policies

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("policies")
}

func resourceURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("policies", id)
}
