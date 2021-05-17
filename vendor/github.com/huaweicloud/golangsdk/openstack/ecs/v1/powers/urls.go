package powers

import "github.com/huaweicloud/golangsdk"

const rootURL = "cloudservers"

func actionURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootURL, "action")
}
