package attachments

import "github.com/chnsz/golangsdk"

func queryURL(client *golangsdk.ServiceClient, instanceId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "attachments")
}
