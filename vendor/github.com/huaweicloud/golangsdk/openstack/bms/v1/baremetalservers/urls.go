package baremetalservers

import "github.com/huaweicloud/golangsdk"

func createURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("baremetalservers")
}

func getURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("baremetalservers", serverID)
}

func jobURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL("jobs", jobId)
}

func putURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("baremetalservers", serverID)
}
