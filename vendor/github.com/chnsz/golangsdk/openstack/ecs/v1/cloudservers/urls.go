package cloudservers

import "github.com/chnsz/golangsdk"

func createURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("cloudservers")
}

func deleteURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("cloudservers", "delete")
}

func getURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("cloudservers", serverID)
}

func resizeURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("cloudservers", serverID, "resize")
}

func listDetailURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("cloudservers", "detail")
}

func jobURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL("jobs", jobId)
}

func orderURL(sc *golangsdk.ServiceClient, orderId string) string {
	return sc.ServiceURL(sc.DomainID, "common/order-mgr/orders-resource", orderId)
}

func deleteOrderURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(sc.DomainID, "common/order-mgr/resources/delete")
}

func passwordURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("cloudservers", id, "os-reset-password")
}

func updateURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("cloudservers", serverID)
}
