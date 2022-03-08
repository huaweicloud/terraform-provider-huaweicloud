package dataset

import "github.com/chnsz/golangsdk"

const (
	dataSetPath = "datasets"
)

// POST /v2/{project_id}/datasets
func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(dataSetPath)
}

// PUT /v2/{project_id}/datasets/{dataset_id}
func updateURL(c *golangsdk.ServiceClient, datasetId string) string {
	return c.ServiceURL(dataSetPath, datasetId)
}

// DELETE /v2/{project_id}/datasets/{dataset_id}
func deleteURL(c *golangsdk.ServiceClient, datasetId string) string {
	return c.ServiceURL(dataSetPath, datasetId)
}

// GET /v2/{project_id}/datasets/{dataset_id}
func getURL(c *golangsdk.ServiceClient, datasetId string) string {
	return c.ServiceURL(dataSetPath, datasetId)
}

// GET /v2/{project_id}/datasets
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(dataSetPath)
}
