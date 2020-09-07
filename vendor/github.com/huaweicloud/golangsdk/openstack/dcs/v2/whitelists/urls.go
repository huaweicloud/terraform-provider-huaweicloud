package whitelists

import "github.com/huaweicloud/golangsdk"

const resourcePath = "instance"

// resourceURL will build the url of put and get request url
// url: client.Endpoint/instance/{instance_id}/whitelist
func resourceURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id, "whitelist")
}
