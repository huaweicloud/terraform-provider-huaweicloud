package eips

import "github.com/chnsz/golangsdk"

const resourcePath = "publicips"

func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("eip", resourcePath, id)
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("eip", resourcePath)
}
