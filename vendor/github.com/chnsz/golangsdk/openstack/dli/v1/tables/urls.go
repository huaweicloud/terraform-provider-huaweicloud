package tables

import "github.com/chnsz/golangsdk"

const (
	databasePath = "databases"
	tablePath    = "tables"
)

// POST /v1.0/{project_id}/databases/{database_name}/tables
func createURL(c *golangsdk.ServiceClient, databaseName string) string {
	return c.ServiceURL(databasePath, databaseName, tablePath)
}

// GET /v1.0/{project_id}/databases/{database_name}/tables
func listURL(c *golangsdk.ServiceClient, databaseName string) string {
	return c.ServiceURL(databasePath, databaseName, tablePath)
}

// DELETE /v1.0/{project_id}/databases/{database_name}/tables/{table_name}
func deleteURL(c *golangsdk.ServiceClient, databaseName string, tableName string) string {
	return c.ServiceURL(databasePath, databaseName, tablePath, tableName)
}

// GET /v1.0/{project_id}/databases/{database_name}/tables/{table_name}
func getURL(c *golangsdk.ServiceClient, databaseName string, tableName string) string {
	return c.ServiceURL(databasePath, databaseName, tablePath, tableName)
}

// GET /v1.0/{project_id}/databases/{database_name}/tables/{table_name}/partitions
func partitionsURL(c *golangsdk.ServiceClient, databaseName string, tableName string) string {
	return c.ServiceURL(databasePath, databaseName, tablePath, tableName, "partitions")
}

// PUT /v1.0/{project_id}/databases/{database_name}/tables/{table_name}/owner
func updateOwnerURL(c *golangsdk.ServiceClient, databaseName string, tableName string) string {
	return c.ServiceURL(databasePath, databaseName, tablePath, tableName, "owner")
}
