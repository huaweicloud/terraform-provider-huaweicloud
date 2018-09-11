package images

import "github.com/huaweicloud/golangsdk"

func listDetailURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("images", "detail")
}

func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}

func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}
