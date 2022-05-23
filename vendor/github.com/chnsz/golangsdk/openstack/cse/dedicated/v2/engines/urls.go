package engines

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("enginemgr", "engines")
}

func resourceURL(client *golangsdk.ServiceClient, engineId string) string {
	return client.ServiceURL("enginemgr", "engines", engineId)
}

func jobURL(client *golangsdk.ServiceClient, engineId, jobId string) string {
	return client.ServiceURL("enginemgr", "engines", engineId, "jobs", jobId)
}
