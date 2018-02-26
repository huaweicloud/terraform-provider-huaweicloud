package regions

import "github.com/huaweicloud/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("regions")
}

func getURL(client *golangsdk.ServiceClient, regionID string) string {
	return client.ServiceURL("regions", regionID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("regions")
}

func updateURL(client *golangsdk.ServiceClient, regionID string) string {
	return client.ServiceURL("regions", regionID)
}

func deleteURL(client *golangsdk.ServiceClient, regionID string) string {
	return client.ServiceURL("regions", regionID)
}
