package services

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("registry", "microservices")
}

func resourceURL(client *golangsdk.ServiceClient, serviceId string) string {
	return client.ServiceURL("registry", "microservices", serviceId)
}
