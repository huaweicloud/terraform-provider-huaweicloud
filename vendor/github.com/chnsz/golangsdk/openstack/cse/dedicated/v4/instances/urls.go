package instances

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, serviceId string) string {
	return client.ServiceURL("registry", "microservices", serviceId, "instances")
}

func resourceURL(client *golangsdk.ServiceClient, serviceId, instanceId string) string {
	return client.ServiceURL("registry", "microservices", serviceId, "instances", instanceId)
}
