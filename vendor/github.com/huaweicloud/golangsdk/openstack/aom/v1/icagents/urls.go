package icagents

import "github.com/huaweicloud/golangsdk"

const (
	rootPath = "agents"
)

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}
