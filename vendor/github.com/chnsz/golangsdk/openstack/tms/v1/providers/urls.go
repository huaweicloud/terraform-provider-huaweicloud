package providers

import "github.com/chnsz/golangsdk"

func queryURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("tms/providers")
}
