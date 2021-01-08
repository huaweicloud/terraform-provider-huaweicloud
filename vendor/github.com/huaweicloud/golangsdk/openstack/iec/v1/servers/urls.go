package servers

import "github.com/huaweicloud/golangsdk"

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("cloudservers")
}

// getURL get iec server detail url
func getURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("cloudservers", serverID)
}

// updateURL get iec server update url
func updateURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("cloudservers", serverID)
}

func deleteAllServersURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("cloudservers", "delete")
}
