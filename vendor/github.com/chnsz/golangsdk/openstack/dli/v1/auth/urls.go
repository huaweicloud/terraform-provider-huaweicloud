package auth

import "github.com/chnsz/golangsdk"

// PUT /v1.0/{project_id}/user-authorization
func grantDataPermissionURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("user-authorization")
}

// GET /v1.0/{project_id}/authorization/privileges
func ListDataPermissionUrl(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("authorization", "privileges")
}

// PUT /v1.0/{project_id}/queues/user-authorization
func grantQueuePermissionURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("queues", "user-authorization")
}

// GET /v1.0/{project_id}/queues/{queue_name}/users
func listQueuePermissionURL(c *golangsdk.ServiceClient, queueName string) string {
	return c.ServiceURL("queues", queueName, "users")
}

// GET /v1.0/{project_id}/databases/{database_name}/users
func listDatabasePermissionURL(c *golangsdk.ServiceClient, databaseName string) string {
	return c.ServiceURL("databases", databaseName, "users")
}

// GET /v1.0/{project_id}/databases/{database_name}/tables/{table_name}/users/{user_name}
func getTablePermissionOfUserURL(c *golangsdk.ServiceClient, databaseName string, tableName string, userName string) string {
	return c.ServiceURL("databases", databaseName, "tables", tableName, "users", userName)
}

// GET /v1.0/{project_id}/databases/{database_name}/tables/{table_name}/users
func listTablePermissionURL(c *golangsdk.ServiceClient, databaseName string, tableName string) string {
	return c.ServiceURL("databases", databaseName, "tables", tableName, "users")
}
