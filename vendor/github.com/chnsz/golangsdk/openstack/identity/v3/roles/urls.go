package roles

import "github.com/chnsz/golangsdk"

const (
	rolePath = "roles"
)

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func getURL(client *golangsdk.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func updateURL(client *golangsdk.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func deleteURL(client *golangsdk.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func listAssignmentsURL(client *golangsdk.ServiceClient, targetType, targetID, actorType, actorID string) string {
	return client.ServiceURL(targetType, targetID, actorType, actorID, rolePath)
}

func assignURL(client *golangsdk.ServiceClient, targetType, targetID, actorType, actorID, roleID string) string {
	return client.ServiceURL(targetType, targetID, actorType, actorID, rolePath, roleID)
}

func assignInheritedURL(client *golangsdk.ServiceClient, domainID, groupID, roleID string) string {
	return client.ServiceURL("OS-INHERIT", "domains", domainID, "groups", groupID,
		"roles", roleID, "inherited_to_projects")
}
