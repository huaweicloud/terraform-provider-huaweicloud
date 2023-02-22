package policyassignments

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, domainId string) string {
	return client.ServiceURL("resource-manager/domains", domainId, "policy-assignments")
}

func resourceURL(client *golangsdk.ServiceClient, domainId, assignmentId string) string {
	return client.ServiceURL("resource-manager/domains", domainId, "policy-assignments", assignmentId)
}

func enableURL(client *golangsdk.ServiceClient, domainId, assignmentId string) string {
	return client.ServiceURL("resource-manager/domains", domainId, "policy-assignments", assignmentId, "enable")
}

func disableURL(client *golangsdk.ServiceClient, domainId, assignmentId string) string {
	return client.ServiceURL("resource-manager/domains", domainId, "policy-assignments", assignmentId, "disable")
}

func queryDefinitionURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("resource-manager/policy-definitions")
}
