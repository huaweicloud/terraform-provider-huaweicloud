package eps_permissions

import "github.com/chnsz/golangsdk"

const rootPath = "OS-PERMISSION"

func userGroupPermissionsURL(client *golangsdk.ServiceClient, enterpriseProjectID, userID, roleID string) string {
	return client.ServiceURL(rootPath, "enterprise-projects", enterpriseProjectID, "groups", userID, "roles", roleID)
}

func userGroupPermissionsGetURL(client *golangsdk.ServiceClient, enterpriseProjectID, userID string) string {
	return client.ServiceURL(rootPath, "enterprise-projects", enterpriseProjectID, "groups", userID, "roles")
}

func userPermissionsURL(client *golangsdk.ServiceClient, enterpriseProjectID, userID, roleID string) string {
	return client.ServiceURL(rootPath, "enterprise-projects", enterpriseProjectID, "users", userID, "roles", roleID)
}

func userPermissionsGetURL(client *golangsdk.ServiceClient, enterpriseProjectID, userID string) string {
	return client.ServiceURL(rootPath, "enterprise-projects", enterpriseProjectID, "users", userID, "roles")
}

func agencyPermissionsURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "subjects/agency/scopes/enterprise-project/role-assignments")
}
