package servers

import "github.com/huaweicloud/golangsdk"

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("servers")
}

func listURL(client *golangsdk.ServiceClient) string {
	return createURL(client)
}

func listDetailURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("servers", "detail")
}

func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id)
}

func getURL(client *golangsdk.ServiceClient, id string) string {
	return deleteURL(client, id)
}

func updateURL(client *golangsdk.ServiceClient, id string) string {
	return deleteURL(client, id)
}

func actionURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}

func metadatumURL(client *golangsdk.ServiceClient, id, key string) string {
	return client.ServiceURL("servers", id, "metadata", key)
}

func metadataURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "metadata")
}

func listAddressesURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "ips")
}

func listAddressesByNetworkURL(client *golangsdk.ServiceClient, id, network string) string {
	return client.ServiceURL("servers", id, "ips", network)
}

func passwordURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "os-server-password")
}
