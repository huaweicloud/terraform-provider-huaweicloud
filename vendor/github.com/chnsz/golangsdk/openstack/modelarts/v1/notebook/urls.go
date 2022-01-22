package notebook

import "github.com/chnsz/golangsdk"

const (
	notebookPath = "notebooks"
	tablePath    = "tables"
)

// POST /v1/{project_id}/notebooks
func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(notebookPath)
}

// GET /v1/{project_id}/notebooks
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(notebookPath)
}

// DELETE /v1/{project_id}/notebooks/{id}
func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id)
}

// GET /v1/{project_id}/notebooks/{id}
func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id)
}

// PUT /v1/{project_id}/notebooks/{id}
func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id)
}

// PATCH /v1/{project_id}/notebooks/{id}/lease
func updateLeaseURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id, "lease")
}

// POST /v1/{project_id}/notebooks/{id}/start
func startURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id, "start")
}

// POST /v1/{project_id}/notebooks/{id}/stop
func stopURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id, "stop")
}

// GET /v1/{project_id}/images
func imagesURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("images")
}

// GET /v1/{project_id}/notebooks/{id}/flavors
func switchableFlavorsURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id, "flavors")
}

// POST /v1/{project_id}/notebooks/{instance_id}/storage
// GET /v1/{project_id}/notebooks/{instance_id}/storage
func mountURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(notebookPath, id, "storage")
}

// GET /v1/{project_id}/notebooks/{instance_id}/storage/{storage_id}
// DELETE /v1/{project_id}/notebooks/{instance_id}/storage/{storage_id}
func mountDetailURL(c *golangsdk.ServiceClient, id string, storageId string) string {
	return c.ServiceURL(notebookPath, id, "storage", storageId)
}
