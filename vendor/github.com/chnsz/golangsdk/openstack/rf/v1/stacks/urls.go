package stacks

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("stacks")
}

func resourceURL(client *golangsdk.ServiceClient, stackName string) string {
	return client.ServiceURL("stacks", stackName)
}

func deploymentURL(client *golangsdk.ServiceClient, stackName string) string {
	return client.ServiceURL("stacks", stackName, "deployments")
}

func eventURL(client *golangsdk.ServiceClient, stackName string) string {
	return client.ServiceURL("stacks", stackName, "events")
}
