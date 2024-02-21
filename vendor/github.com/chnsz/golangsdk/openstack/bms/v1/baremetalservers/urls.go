package baremetalservers

import "github.com/chnsz/golangsdk"

// /baremetalservers/
const resourcePath = "baremetalservers"

func createURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(resourcePath)
}

func getURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL(resourcePath, serverID)
}

func listURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(resourcePath, "detail")
}

func jobURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL("jobs", jobId)
}

func putURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL(resourcePath, serverID)
}

func deleteNicsURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL(resourcePath, serverID, "nics", "delete")
}

func addNicsURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL(resourcePath, serverID, "nics")
}

func serverStatusPostURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(resourcePath, "action")
}

func metadataURL(client *golangsdk.ServiceClient, serverID string) string {
	return client.ServiceURL("baremetalservers", serverID, "metadata")
}
