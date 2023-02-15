package interfaces

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("dcaas/virtual-interfaces")
}

func resourceURL(client *golangsdk.ServiceClient, interfaceId string) string {
	return client.ServiceURL("dcaas/virtual-interfaces", interfaceId)
}
