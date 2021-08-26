package dependencies

import "github.com/huaweicloud/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("fgs", "dependencies")
}
