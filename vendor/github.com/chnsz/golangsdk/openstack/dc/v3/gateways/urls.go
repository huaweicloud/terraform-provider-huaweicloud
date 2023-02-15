package gateways

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("dcaas/virtual-gateways")
}

func resourceURL(client *golangsdk.ServiceClient, gatewayId string) string {
	return client.ServiceURL("dcaas/virtual-gateways", gatewayId)
}
