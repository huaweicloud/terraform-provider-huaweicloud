package version

import "github.com/chnsz/golangsdk"

const (
	dataSetPath  = "datasets"
	versionsPath = "versions"
)

// POST /v2/{project_id}/datasets/{dataset_id}/versions
func createURL(c *golangsdk.ServiceClient, datasetId string) string {
	return c.ServiceURL(dataSetPath, datasetId, versionsPath)
}

// DELETE /v2/{project_id}/datasets/{dataset_id}/versions/{version_id}
func deleteURL(c *golangsdk.ServiceClient, datasetId string, versionId string) string {
	return c.ServiceURL(dataSetPath, datasetId, versionsPath, versionId)
}

// GET /v2/{project_id}/datasets/{dataset_id}/versions/{version_id}
func getURL(c *golangsdk.ServiceClient, datasetId string, versionId string) string {
	return c.ServiceURL(dataSetPath, datasetId, versionsPath, versionId)
}

// GET /v2/{project_id}/datasets/{dataset_id}/versions
func listURL(c *golangsdk.ServiceClient, datasetId string) string {
	return c.ServiceURL(dataSetPath, datasetId, versionsPath)
}
