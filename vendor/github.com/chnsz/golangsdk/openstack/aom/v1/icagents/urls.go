package icagents

import "github.com/chnsz/golangsdk"

const (
	rootPath = "agents"
)

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}
