package catalog

import "github.com/huaweicloud/golangsdk"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("auth/catalog")
}
