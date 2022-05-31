package instances

import "github.com/chnsz/golangsdk"

const rootPath = "cas/applications"

func rootURL(client *golangsdk.ServiceClient, appId, componentId string) string {
	return client.ServiceURL(rootPath, appId, "components", componentId, "instances")
}

func resourceURL(client *golangsdk.ServiceClient, appId, componentId, instanceId string) string {
	return client.ServiceURL(rootPath, appId, "components", componentId, "instances", instanceId)
}

func actionURL(client *golangsdk.ServiceClient, appId, componentId, instanceId string) string {
	return client.ServiceURL(rootPath, appId, "components", componentId, "instances", instanceId, "action")
}
