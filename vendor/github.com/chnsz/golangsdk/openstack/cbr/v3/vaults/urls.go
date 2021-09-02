package vaults

import "github.com/chnsz/golangsdk"

const resourcePath = "vaults"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

func resourceURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

func addResourcesURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id, "addresources")
}

func removeResourcesURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id, "removeresources")
}

func migrateResourcesURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id, "migrateresources")
}

func bindPolicyURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id, "associatepolicy")
}

func unbindPolicyURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id, "dissociatepolicy")
}
